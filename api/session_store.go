package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type SessionStatus string

const (
	SessionQueued      SessionStatus = "queued"
	SessionPreparing   SessionStatus = "preparing"
	SessionRunning     SessionStatus = "running"
	SessionCompleted   SessionStatus = "completed"
	SessionFailed      SessionStatus = "failed"
	SessionInterrupted SessionStatus = "interrupted"
	SessionCancelled   SessionStatus = "cancelled"
)

type Session struct {
	ID              string        `json:"id"`
	StoryID         string        `json:"storyId"`
	Status          SessionStatus `json:"status"`
	Kind            SessionKind   `json:"kind"`
	Worktree        string        `json:"worktree"`
	ThreadID        string        `json:"threadId"`
	Prompt          string        `json:"prompt"`
	Error           string        `json:"error"`
	PlanSubmitted   bool          `json:"planSubmitted"`
	CreatedAt       string        `json:"createdAt"`
	FinishedAt      string        `json:"finishedAt"`
	ParentSessionID string        `json:"parentSessionId"`
	FeedbackBatchID string        `json:"feedbackBatchId"`
	ReviewStatus    ReviewStatus  `json:"reviewStatus"`
	ReviewedAt      string        `json:"reviewedAt"`
}
type SessionEvent struct {
	ID     int64  `json:"id"`
	Stream string `json:"stream"`
	Data   []byte `json:"data"`
}
type ActionItem struct {
	ID                 string   `json:"id"`
	SessionID          string   `json:"sessionId"`
	Key                string   `json:"key"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	AcceptanceCriteria []string `json:"acceptanceCriteria"`
}

const sessionColumns = "id, story_id, status, worktree, thread_id, prompt, error, plan_submitted, created_at, finished_at, COALESCE(parent_session_id,''), COALESCE(feedback_batch_id,''), review_status, reviewed_at, kind"

func migrateSessions(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`CREATE TABLE sessions (
 id TEXT PRIMARY KEY NOT NULL,
 story_id TEXT NOT NULL REFERENCES user_stories(id),
 status TEXT NOT NULL CHECK (status IN ('queued','preparing','running','completed','failed','interrupted','cancelled')),
 worktree TEXT NOT NULL DEFAULT '', thread_id TEXT NOT NULL DEFAULT '', prompt TEXT NOT NULL DEFAULT '', error TEXT NOT NULL DEFAULT '',
 plan_submitted INTEGER NOT NULL DEFAULT 0,
 token_hash BLOB, process_id INTEGER NOT NULL DEFAULT 0, process_identity TEXT NOT NULL DEFAULT '',
 created_at TEXT NOT NULL, finished_at TEXT NOT NULL DEFAULT ''
 );
 CREATE UNIQUE INDEX one_active_session_per_story ON sessions(story_id) WHERE status IN ('queued','preparing','running');
 CREATE INDEX sessions_by_story ON sessions(story_id);
 CREATE TABLE session_events (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 session_id TEXT NOT NULL REFERENCES sessions(id),
 stream TEXT NOT NULL, data BLOB NOT NULL,
 created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
 );
 CREATE INDEX events_by_session ON session_events(session_id,id);
 CREATE TABLE action_items (
 id TEXT PRIMARY KEY NOT NULL,
 session_id TEXT NOT NULL REFERENCES sessions(id),
 item_key TEXT NOT NULL, title TEXT NOT NULL, description TEXT NOT NULL,
 acceptance_criteria TEXT NOT NULL,
 UNIQUE(session_id,item_key)
 );
 PRAGMA user_version = 4;`)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return migratePlanReview(db)
}
func scanSession(row interface{ Scan(...any) error }) (Session, error) {
	var session Session
	err := row.Scan(&session.ID, &session.StoryID, &session.Status, &session.Worktree, &session.ThreadID, &session.Prompt, &session.Error, &session.PlanSubmitted, &session.CreatedAt, &session.FinishedAt, &session.ParentSessionID, &session.FeedbackBatchID, &session.ReviewStatus, &session.ReviewedAt, &session.Kind)
	return session, err
}
func (s *Store) session(id string) (Session, error) {
	return scanSession(s.db.QueryRow("SELECT "+sessionColumns+" FROM sessions WHERE id = ?", id))
}

// Failed/interrupted attempts are deliberately not auto-retried: no surprise retry spending.
func (s *Store) enqueueStories() error {
	_, err := s.db.Exec(`INSERT INTO sessions (id, story_id, status, created_at)
 SELECT lower(hex(randomblob(16))), id, 'queued', strftime('%Y-%m-%dT%H:%M:%fZ','now') FROM user_stories
 WHERE length(trim(goal)) > 0 AND NOT EXISTS (SELECT 1 FROM sessions WHERE story_id = user_stories.id)`)
	return err
}
func (s *Store) claimSession() (Session, error) {
	return scanSession(s.db.QueryRow(`UPDATE sessions SET status = 'preparing'
 WHERE id = (SELECT id FROM sessions WHERE status = 'queued' ORDER BY rowid LIMIT 1)
 RETURNING ` + sessionColumns))
}
func (s *Store) retryStory(storyID string) (Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Session{}, err
	}
	defer tx.Rollback()
	var goal string
	if err = tx.QueryRow("SELECT goal FROM user_stories WHERE id = ?", storyID).Scan(&goal); err != nil {
		return Session{}, err
	}
	if strings.TrimSpace(goal) == "" {
		return Session{}, fmt.Errorf("User story needs a goal before planning")
	}
	var active int
	if err = tx.QueryRow("SELECT count(*) FROM sessions WHERE story_id = ? AND status IN ('queued','preparing','running')", storyID).Scan(&active); err != nil {
		return Session{}, err
	}
	if active > 0 {
		return Session{}, fmt.Errorf("User story already has an active session")
	}
	var approvedCount int
	if err := tx.QueryRow("SELECT count(*) FROM sessions WHERE story_id=? AND review_status='approved' AND kind='planning'", storyID).Scan(&approvedCount); err != nil {
		return Session{}, err
	}
	if approvedCount > 0 {
		return Session{}, fmt.Errorf("%w: approved stories must retry implementation, not planning", ErrReviewConflict)
	}
	var latestID, parentID, batchID string
	var review ReviewStatus
	err = tx.QueryRow("SELECT id,COALESCE(parent_session_id,''),COALESCE(feedback_batch_id,''),review_status FROM sessions WHERE story_id=? ORDER BY rowid DESC LIMIT 1", storyID).Scan(&latestID, &parentID, &batchID, &review)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Session{}, err
	}
	if review == ReviewApproved {
		return Session{}, fmt.Errorf("%w: approved plans cannot be replanned", ErrReviewConflict)
	}
	if latestID != "" {
		var pending int
		if err := tx.QueryRow("SELECT count(*) FROM plan_comments WHERE session_id=? AND batch_id IS NULL", latestID).Scan(&pending); err != nil {
			return Session{}, err
		}
		if pending > 0 {
			return Session{}, fmt.Errorf("%w: send pending feedback instead of replanning", ErrReviewConflict)
		}
	}
	if batchID != "" {
		session, err := queueRevision(tx, storyID, parentID, batchID)
		if err != nil {
			return Session{}, err
		}
		if err := tx.Commit(); err != nil {
			return Session{}, err
		}
		return session, nil
	}
	session := Session{ID: rand.Text(), StoryID: storyID, Status: SessionQueued, Kind: PlanningSession, ReviewStatus: ReviewPending, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if _, err = tx.Exec("INSERT INTO sessions (id,story_id,status,created_at) VALUES (?,?,?,?)", session.ID, storyID, session.Status, session.CreatedAt); err != nil {
		return Session{}, err
	}
	if err = tx.Commit(); err != nil {
		return Session{}, err
	}
	return session, nil
}
func (s *Store) appendEvent(sessionID, stream string, data []byte) error {
	_, err := s.db.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,?,?)", sessionID, stream, data)
	return err
}
func (s *Store) events(sessionID string, after int64) ([]SessionEvent, error) {
	rows, err := s.db.Query("SELECT id,stream,data FROM session_events WHERE session_id = ? AND id > ? ORDER BY id LIMIT 200", sessionID, after)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := []SessionEvent{}
	for rows.Next() {
		var event SessionEvent
		if err := rows.Scan(&event.ID, &event.Stream, &event.Data); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
func (s *Store) items(sessionID string) ([]ActionItem, error) {
	rows, err := s.db.Query("SELECT id,session_id,item_key,title,description,acceptance_criteria FROM action_items WHERE session_id = ? ORDER BY rowid", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ActionItem{}
	for rows.Next() {
		var item ActionItem
		var criteria string
		if err := rows.Scan(&item.ID, &item.SessionID, &item.Key, &item.Title, &item.Description, &criteria); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(criteria), &item.AcceptanceCriteria); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (s *Store) startSession(id, worktree, prompt, token string) error {
	hash := sha256.Sum256([]byte(token))
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec("UPDATE sessions SET status='running',worktree=?,prompt=?,token_hash=? WHERE id=? AND status='preparing'", worktree, prompt, hash[:], id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err == nil && n != 1 {
		return fmt.Errorf("session is not preparing")
	}
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE work_items SET status='working' WHERE current_session_id=? AND status='queued'", id); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) finishSession(id string, status SessionStatus, message string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec("UPDATE sessions SET status=?,error=?,finished_at=?,token_hash=NULL WHERE id=? AND (status IN ('preparing','running') OR (status='queued' AND ?='cancelled'))", status, message, time.Now().UTC().Format(time.RFC3339Nano), id, status)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	if err := finishWorkItemSession(tx, id, status, message); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) authenticatedSession(token string) (Session, error) {
	if token == "" {
		return Session{}, sql.ErrNoRows
	}
	hash := sha256.Sum256([]byte(token))
	return scanSession(s.db.QueryRow("SELECT "+sessionColumns+" FROM sessions WHERE token_hash=? AND status='running'", hash[:]))
}
func (s *Store) storyContext(storyID string) (factoryID, title, goal string, err error) {
	err = s.db.QueryRow("SELECT factory_id,title,goal FROM user_stories WHERE id=?", storyID).Scan(&factoryID, &title, &goal)
	return
}

var ErrPlanClosed = errors.New("session is not accepting plan items")
