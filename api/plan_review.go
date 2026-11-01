package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

type ReviewStatus string

const (
	ReviewPending          ReviewStatus = "pending"
	ReviewApproved         ReviewStatus = "approved"
	ReviewChangesRequested ReviewStatus = "changes_requested"
)

var ErrReviewConflict = errors.New("plan is not available for review")

type PlanComment struct {
	ID             string `json:"id"`
	SessionID      string `json:"sessionId"`
	ItemID         string `json:"itemId"`
	Field          string `json:"field"`
	CriterionIndex int    `json:"criterionIndex"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
	Quote          string `json:"quote"`
	Body           string `json:"body"`
	BatchID        string `json:"batchId"`
	CreatedAt      string `json:"createdAt"`
}
type CommentInput struct {
	ItemID         string `json:"itemId"`
	Field          string `json:"field"`
	CriterionIndex int    `json:"criterionIndex"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
	Quote          string `json:"quote"`
	Body           string `json:"body"`
}

const commentColumns = "id,session_id,item_id,field,criterion_index,start_offset,end_offset,quote,body,COALESCE(batch_id,''),created_at"

func migratePlanReview(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`CREATE TABLE plan_feedback_batches (
 id TEXT PRIMARY KEY NOT NULL, session_id TEXT NOT NULL UNIQUE REFERENCES sessions(id), created_at TEXT NOT NULL
 );
 ALTER TABLE sessions ADD COLUMN parent_session_id TEXT REFERENCES sessions(id);
 ALTER TABLE sessions ADD COLUMN feedback_batch_id TEXT REFERENCES plan_feedback_batches(id);
 ALTER TABLE sessions ADD COLUMN review_status TEXT NOT NULL DEFAULT 'pending' CHECK (review_status IN ('pending','approved','changes_requested'));
 ALTER TABLE sessions ADD COLUMN reviewed_at TEXT NOT NULL DEFAULT '';
 CREATE TABLE plan_comments (
 id TEXT PRIMARY KEY NOT NULL,session_id TEXT NOT NULL REFERENCES sessions(id),item_id TEXT NOT NULL REFERENCES action_items(id),
 field TEXT NOT NULL CHECK (field IN ('title','description','acceptanceCriteria')),criterion_index INTEGER NOT NULL,
 start_offset INTEGER NOT NULL CHECK(start_offset>=0),end_offset INTEGER NOT NULL CHECK(end_offset>start_offset),
 quote TEXT NOT NULL,body TEXT NOT NULL,batch_id TEXT REFERENCES plan_feedback_batches(id),created_at TEXT NOT NULL
 );
 CREATE INDEX comments_by_plan ON plan_comments(session_id);
 CREATE INDEX comments_by_batch ON plan_comments(batch_id);
 CREATE TABLE feedback_resolutions (
 session_id TEXT NOT NULL REFERENCES sessions(id),comment_id TEXT NOT NULL REFERENCES plan_comments(id),resolution TEXT NOT NULL,
 PRIMARY KEY(session_id,comment_id)
 );
 PRAGMA user_version = 5;`)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return migrateImplementations(db)
}
func reviewable(tx *sql.Tx, id string) (Session, error) {
	session, err := scanSession(tx.QueryRow("SELECT "+sessionColumns+" FROM sessions WHERE id=?", id))
	if err != nil {
		return Session{}, err
	}
	var latest string
	if err := tx.QueryRow("SELECT id FROM sessions WHERE story_id=? ORDER BY rowid DESC LIMIT 1", session.StoryID).Scan(&latest); err != nil {
		return Session{}, err
	}
	if latest != id {
		return Session{}, fmt.Errorf("%w: only the latest plan can be reviewed", ErrReviewConflict)
	}
	if session.Status != SessionCompleted || !session.PlanSubmitted || session.ReviewStatus != ReviewPending {
		return Session{}, ErrReviewConflict
	}
	return session, nil
}
func validComment(input CommentInput, text string) error {
	if input.Start < 0 || input.End <= input.Start || !utf8.ValidString(input.Quote) || !utf8.ValidString(input.Body) || strings.TrimSpace(input.Body) == "" || utf8.RuneCountInString(input.Body) > 4000 {
		return fmt.Errorf("invalid feedback selection or comment")
	}
	units := utf16.Encode([]rune(text))
	quote := utf16.Encode([]rune(input.Quote))
	if input.End > len(units) || input.End-input.Start != len(quote) {
		return fmt.Errorf("selected text no longer matches the plan")
	}
	for index, unit := range quote {
		if units[input.Start+index] != unit {
			return fmt.Errorf("selected text no longer matches the plan")
		}
	}
	if len(quote) == 0 {
		return fmt.Errorf("select some plan text")
	}
	return nil
}
func (s *Store) addPlanComment(sessionID string, input CommentInput) (PlanComment, error) {
	input.Body = strings.TrimSpace(input.Body)
	tx, err := s.db.Begin()
	if err != nil {
		return PlanComment{}, err
	}
	defer tx.Rollback()
	if _, err := reviewable(tx, sessionID); err != nil {
		return PlanComment{}, err
	}
	var title, description, criteriaJSON string
	if err := tx.QueryRow("SELECT title,description,acceptance_criteria FROM action_items WHERE id=? AND session_id=?", input.ItemID, sessionID).Scan(&title, &description, &criteriaJSON); err != nil {
		return PlanComment{}, err
	}
	var text string
	switch input.Field {
	case "title":
		if input.CriterionIndex != -1 {
			return PlanComment{}, fmt.Errorf("invalid item selection")
		}
		text = title
	case "description":
		if input.CriterionIndex != -1 {
			return PlanComment{}, fmt.Errorf("invalid work selection")
		}
		text = description
	case "acceptanceCriteria":
		var criteria []string
		if err := json.Unmarshal([]byte(criteriaJSON), &criteria); err != nil {
			return PlanComment{}, err
		}
		if input.CriterionIndex < 0 || input.CriterionIndex >= len(criteria) {
			return PlanComment{}, fmt.Errorf("invalid acceptance criterion")
		}
		text = criteria[input.CriterionIndex]
	default:
		return PlanComment{}, fmt.Errorf("invalid plan field")
	}
	if err := validComment(input, text); err != nil {
		return PlanComment{}, err
	}
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM plan_comments WHERE session_id=? AND batch_id IS NULL", sessionID).Scan(&count); err != nil {
		return PlanComment{}, err
	}
	if count >= 100 {
		return PlanComment{}, fmt.Errorf("feedback limit reached (100 comments)")
	}
	comment := PlanComment{ID: rand.Text(), SessionID: sessionID, ItemID: input.ItemID, Field: input.Field, CriterionIndex: input.CriterionIndex, Start: input.Start, End: input.End, Quote: input.Quote, Body: input.Body, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	_, err = tx.Exec("INSERT INTO plan_comments (id,session_id,item_id,field,criterion_index,start_offset,end_offset,quote,body,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)", comment.ID, sessionID, comment.ItemID, comment.Field, comment.CriterionIndex, comment.Start, comment.End, comment.Quote, comment.Body, comment.CreatedAt)
	if err != nil {
		return PlanComment{}, err
	}
	if err := tx.Commit(); err != nil {
		return PlanComment{}, err
	}
	return comment, nil
}
func (s *Store) removePlanComment(sessionID, id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := reviewable(tx, sessionID); err != nil {
		return err
	}
	result, err := tx.Exec("DELETE FROM plan_comments WHERE id=? AND session_id=? AND batch_id IS NULL", id, sessionID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}
func scanComments(rows *sql.Rows) ([]PlanComment, error) {
	defer rows.Close()
	comments := []PlanComment{}
	for rows.Next() {
		var comment PlanComment
		if err := rows.Scan(&comment.ID, &comment.SessionID, &comment.ItemID, &comment.Field, &comment.CriterionIndex, &comment.Start, &comment.End, &comment.Quote, &comment.Body, &comment.BatchID, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}
func (s *Store) planComments(sessionID string) ([]PlanComment, error) {
	rows, err := s.db.Query("SELECT "+commentColumns+" FROM plan_comments WHERE session_id=? ORDER BY rowid", sessionID)
	if err != nil {
		return nil, err
	}
	return scanComments(rows)
}
func (s *Store) receivedFeedback(batchID string) ([]PlanComment, error) {
	rows, err := s.db.Query("SELECT "+commentColumns+" FROM plan_comments WHERE batch_id=? ORDER BY rowid", batchID)
	if err != nil {
		return nil, err
	}
	return scanComments(rows)
}
func (s *Store) approvePlan(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := reviewable(tx, id); err != nil {
		return err
	}
	var pending, items int
	if err := tx.QueryRow("SELECT count(*) FROM plan_comments WHERE session_id=? AND batch_id IS NULL", id).Scan(&pending); err != nil {
		return err
	}
	if pending > 0 {
		return fmt.Errorf("%w: send or remove pending feedback before approving", ErrReviewConflict)
	}
	if err := tx.QueryRow("SELECT count(*) FROM action_items WHERE session_id=?", id).Scan(&items); err != nil {
		return err
	}
	if items == 0 {
		return ErrReviewConflict
	}
	if _, err := tx.Exec("UPDATE sessions SET review_status='approved',reviewed_at=? WHERE id=?", time.Now().UTC().Format(time.RFC3339Nano), id); err != nil {
		return err
	}
	if _, err := queueImplementation(tx, id); err != nil {
		return err
	}
	return tx.Commit()
}
func clonePlan(tx *sql.Tx, sourceID, targetID string) error {
	_, err := tx.Exec(`INSERT INTO action_items (id,session_id,item_key,title,description,acceptance_criteria)
 SELECT lower(hex(randomblob(16))),?,item_key,title,description,acceptance_criteria FROM action_items WHERE session_id=? ORDER BY rowid`, targetID, sourceID)
	return err
}
func queueRevision(tx *sql.Tx, storyID, parentID, batchID string) (Session, error) {
	session := Session{ID: rand.Text(), StoryID: storyID, Status: SessionQueued, Kind: PlanningSession, ReviewStatus: ReviewPending, ParentSessionID: parentID, FeedbackBatchID: batchID, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	_, err := tx.Exec("INSERT INTO sessions (id,story_id,status,parent_session_id,feedback_batch_id,created_at) VALUES (?,?,'queued',?,?,?)", session.ID, storyID, parentID, batchID, session.CreatedAt)
	if err != nil {
		return Session{}, err
	}
	if err := clonePlan(tx, parentID, session.ID); err != nil {
		return Session{}, err
	}
	return session, nil
}
func (s *Store) sendPlanFeedback(id string) (Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Session{}, err
	}
	defer tx.Rollback()
	// Retrying a delivered batch must not create another Codex process.
	var review ReviewStatus
	if err := tx.QueryRow("SELECT review_status FROM sessions WHERE id=?", id).Scan(&review); err != nil {
		return Session{}, err
	}
	if review == ReviewChangesRequested {
		return scanSession(tx.QueryRow("SELECT "+sessionColumns+" FROM sessions WHERE parent_session_id=? ORDER BY rowid DESC LIMIT 1", id))
	}
	source, err := reviewable(tx, id)
	if err != nil {
		return Session{}, err
	}
	var pending int
	if err := tx.QueryRow("SELECT count(*) FROM plan_comments WHERE session_id=? AND batch_id IS NULL", id).Scan(&pending); err != nil {
		return Session{}, err
	}
	if pending == 0 {
		return Session{}, fmt.Errorf("%w: add feedback before sending", ErrReviewConflict)
	}
	batchID := rand.Text()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := tx.Exec("INSERT INTO plan_feedback_batches (id,session_id,created_at) VALUES (?,?,?)", batchID, id, now); err != nil {
		return Session{}, err
	}
	if _, err := tx.Exec("UPDATE plan_comments SET batch_id=? WHERE session_id=? AND batch_id IS NULL", batchID, id); err != nil {
		return Session{}, err
	}
	if _, err := tx.Exec("UPDATE sessions SET review_status='changes_requested',reviewed_at=? WHERE id=?", now, id); err != nil {
		return Session{}, err
	}
	next, err := queueRevision(tx, source.StoryID, id, batchID)
	if err != nil {
		return Session{}, err
	}
	event, _ := json.Marshal(map[string]any{"type": "feedback_batch", "sourceSessionId": id, "batchId": batchID, "comments": pending})
	if _, err := tx.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,'review',?)", next.ID, event); err != nil {
		return Session{}, err
	}
	if err := tx.Commit(); err != nil {
		return Session{}, err
	}
	return next, nil
}
