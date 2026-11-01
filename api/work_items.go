package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type WorkItem struct {
	ActionItem
	PlanSessionID     string                 `json:"planSessionId"`
	Position          int                    `json:"position"`
	Status            string                 `json:"status"`
	Blocker           string                 `json:"blocker"`
	CurrentSessionID  string                 `json:"currentSessionId"`
	LastGoalSessionID string                 `json:"lastGoalSessionId"`
	Criteria          []AcceptanceSubmission `json:"criteria"`
	Reviews           []ExecutiveReview      `json:"reviews"`
	Replies           []string               `json:"replies"`
}
type AcceptanceSubmission struct {
	CriterionIndex int          `json:"criterionIndex"`
	Summary        string       `json:"summary"`
	Evidence       []string     `json:"evidence"`
	ArtifactTree   ArtifactTree `json:"artifactTree"`
	CheckIDs       []string     `json:"checkIds"`
}
type ExecutiveReview struct {
	ID               string       `json:"id"`
	WorkItemID       string       `json:"workItemId"`
	GoalSessionID    string       `json:"goalSessionId"`
	SummarySessionID string       `json:"summarySessionId"`
	Summary          string       `json:"summary"`
	Verification     string       `json:"verification"`
	Risks            []string     `json:"risks"`
	Decision         ReviewStatus `json:"decision"`
	Feedback         string       `json:"feedback"`
	CreatedAt        string       `json:"createdAt"`
	ReviewedAt       string       `json:"reviewedAt"`
	ArtifactTree     ArtifactTree `json:"artifactTree"`
	ApprovedCommit   string       `json:"approvedCommit"`
}
type GoalWorkflow struct {
	PlanSessionID string       `json:"planSessionId"`
	Items         []WorkItem   `json:"items"`
	Integration   *Integration `json:"integration"`
}

const workItemColumns = `item.id,item.session_id,item.item_key,item.title,item.description,item.acceptance_criteria,work.plan_session_id,work.position,work.status,work.blocker,COALESCE(work.current_session_id,''),COALESCE(work.last_goal_session_id,'')`

func scanWorkItem(row interface{ Scan(...any) error }) (WorkItem, error) {
	var item WorkItem
	var criteria string
	err := row.Scan(&item.ID, &item.SessionID, &item.Key, &item.Title, &item.Description, &criteria, &item.PlanSessionID, &item.Position, &item.Status, &item.Blocker, &item.CurrentSessionID, &item.LastGoalSessionID)
	if err == nil {
		err = json.Unmarshal([]byte(criteria), &item.AcceptanceCriteria)
	}
	item.Replies = []string{}
	item.Criteria = []AcceptanceSubmission{}
	item.Reviews = []ExecutiveReview{}
	return item, err
}
func (s *Store) workItem(id string) (WorkItem, error) {
	item, err := scanWorkItem(s.db.QueryRow("SELECT "+workItemColumns+" FROM work_items AS work JOIN action_items AS item ON item.id=work.id WHERE work.id=?", id))
	if err != nil {
		return item, err
	}
	rows, err := s.db.Query("SELECT criterion_index,summary,evidence,artifact_tree,check_ids FROM acceptance_submissions WHERE work_item_id=? AND session_id=? ORDER BY criterion_index", id, item.LastGoalSessionID)
	if err != nil {
		return item, err
	}
	for rows.Next() {
		var submission AcceptanceSubmission
		var evidence, ids string
		if err := rows.Scan(&submission.CriterionIndex, &submission.Summary, &evidence, &submission.ArtifactTree, &ids); err != nil {
			_ = rows.Close()
			return item, err
		}
		if err := json.Unmarshal([]byte(evidence), &submission.Evidence); err != nil {
			_ = rows.Close()
			return item, err
		}
		if err := json.Unmarshal([]byte(ids), &submission.CheckIDs); err != nil {
			_ = rows.Close()
			return item, err
		}
		item.Criteria = append(item.Criteria, submission)
	}
	err = errors.Join(rows.Err(), rows.Close())
	if err != nil {
		return item, err
	}
	item.Reviews, err = s.workItemReviews(id)
	if err != nil {
		return item, err
	}
	messages, err := s.db.Query("SELECT body FROM work_item_messages WHERE work_item_id=? ORDER BY id", id)
	if err != nil {
		return item, err
	}
	defer messages.Close()
	for messages.Next() {
		var body string
		if err := messages.Scan(&body); err != nil {
			return item, err
		}
		item.Replies = append(item.Replies, body)
	}
	return item, messages.Err()
}
func (s *Store) workflow(planID string) (*GoalWorkflow, error) {
	var exists int
	if err := s.db.QueryRow("SELECT 1 FROM goal_workflows WHERE plan_session_id=?", planID).Scan(&exists); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id FROM work_items WHERE plan_session_id=? ORDER BY position", planID)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return nil, err
	}
	result := &GoalWorkflow{PlanSessionID: planID, Items: []WorkItem{}}
	for _, id := range ids {
		item, err := s.workItem(id)
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, item)
	}
	result.Integration, err = s.goalIntegration(planID)
	return result, err
}
func queueWorkItemSession(tx *sql.Tx, itemID string, kind SessionKind) (Session, error) {
	var planID, storyID string
	if err := tx.QueryRow("SELECT work.plan_session_id,plan.story_id FROM work_items AS work JOIN sessions AS plan ON plan.id=work.plan_session_id WHERE work.id=? AND plan.review_status='approved'", itemID).Scan(&planID, &storyID); err != nil {
		return Session{}, err
	}
	session := Session{ID: rand.Text(), StoryID: storyID, Status: SessionQueued, Kind: kind, ParentSessionID: planID, ReviewStatus: ReviewPending, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if _, err := tx.Exec("INSERT INTO sessions (id,story_id,status,kind,parent_session_id,created_at) VALUES (?,?,'queued',?,?,?)", session.ID, storyID, kind, planID, session.CreatedAt); err != nil {
		return Session{}, err
	}
	if _, err := tx.Exec("INSERT INTO work_item_sessions (session_id,work_item_id) VALUES (?,?)", session.ID, itemID); err != nil {
		return Session{}, err
	}
	state := "summarizing"
	if kind == ImplementationSession {
		state = "queued"
		if _, err := tx.Exec("UPDATE work_items SET last_goal_session_id=?,blocker='' WHERE id=?", session.ID, itemID); err != nil {
			return Session{}, err
		}
	}
	if _, err := tx.Exec("UPDATE work_items SET status=?,current_session_id=? WHERE id=?", state, session.ID, itemID); err != nil {
		return Session{}, err
	}
	event, _ := json.Marshal(map[string]any{"type": "work_item_queued", "workItemId": itemID, "kind": kind})
	if _, err := tx.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,'workflow',?)", session.ID, event); err != nil {
		return Session{}, err
	}
	return session, nil
}
func queueFirstWorkItem(tx *sql.Tx, planID string) (Session, error) {
	var status SessionStatus
	var review ReviewStatus
	var kind SessionKind
	var sealed bool
	if err := tx.QueryRow("SELECT status,review_status,kind,plan_submitted FROM sessions WHERE id=?", planID).Scan(&status, &review, &kind, &sealed); err != nil {
		return Session{}, err
	}
	if status != SessionCompleted || review != ReviewApproved || kind != PlanningSession || !sealed {
		return Session{}, fmt.Errorf("work requires an approved, submitted plan")
	}
	if _, err := tx.Exec("INSERT INTO goal_workflows(plan_session_id) VALUES (?)", planID); err != nil {
		return Session{}, err
	}
	if _, err := tx.Exec(`INSERT INTO work_items (id,plan_session_id,position) SELECT id,session_id,(SELECT count(*)-1 FROM action_items WHERE session_id=item.session_id AND rowid<=item.rowid) FROM action_items AS item WHERE session_id=?`, planID); err != nil {
		return Session{}, err
	}
	var id string
	if err := tx.QueryRow("SELECT id FROM work_items WHERE plan_session_id=? ORDER BY position LIMIT 1", planID).Scan(&id); err != nil {
		return Session{}, err
	}
	return queueWorkItemSession(tx, id, ImplementationSession)
}
func (s *Store) queueApprovedPlans() error {
	if err := s.processIntegrations(); err != nil {
		return err
	}
	rows, err := s.db.Query(`SELECT id FROM sessions AS plan WHERE kind='planning' AND review_status='approved'
 AND NOT EXISTS(SELECT 1 FROM goal_workflows WHERE plan_session_id=plan.id)
 AND NOT EXISTS(SELECT 1 FROM sessions WHERE story_id=plan.story_id AND status IN ('queued','preparing','running'))`)
	if err != nil {
		return err
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return err
	}
	for _, id := range ids {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		var exists int
		err = tx.QueryRow("SELECT count(*) FROM goal_workflows WHERE plan_session_id=?", id).Scan(&exists)
		if err == nil && exists == 0 {
			_, err = queueFirstWorkItem(tx, id)
		}
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return s.queueRequestedSummaries()
}
func (s *Store) queueRequestedSummaries() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var itemID string
	err = tx.QueryRow(`SELECT work.id FROM work_items AS work JOIN sessions AS run ON run.id=work.current_session_id
 WHERE work.status='review_requested' AND run.status='completed'
 AND NOT EXISTS(SELECT 1 FROM sessions WHERE story_id=run.story_id AND status IN ('queued','preparing','running')) ORDER BY run.rowid LIMIT 1`).Scan(&itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	if _, err := queueWorkItemSession(tx, itemID, SummarySession); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) latestImplementation(planID string) (Session, error) {
	return scanSession(s.db.QueryRow("SELECT "+sessionColumns+" FROM sessions WHERE parent_session_id=? AND kind IN ('implementation','summary') ORDER BY rowid DESC LIMIT 1", planID))
}
func (s *Store) retryImplementation(storyID string) (Session, error) {
	var previousID string
	var previousKind SessionKind
	var sealed ArtifactTree
	err := s.db.QueryRow("SELECT COALESCE(work.current_session_id,''),COALESCE(run.kind,'implementation'),work.review_tree FROM work_items AS work JOIN sessions AS plan ON plan.id=work.plan_session_id LEFT JOIN sessions AS run ON run.id=work.current_session_id WHERE plan.story_id=? AND work.status<>'approved' ORDER BY work.position LIMIT 1", storyID).Scan(&previousID, &previousKind, &sealed)
	if err != nil {
		return Session{}, err
	}
	summaryCurrent := false
	if previousKind == SummarySession && sealed != "" {
		_, tree, err := s.sessionArtifact(context.Background(), previousID)
		if err != nil {
			return Session{}, err
		}
		summaryCurrent = tree == sealed
	}
	tx, err := s.db.Begin()
	if err != nil {
		return Session{}, err
	}
	defer tx.Rollback()
	var itemID, state string
	err = tx.QueryRow(`SELECT work.id,work.status FROM work_items AS work JOIN sessions AS plan ON plan.id=work.plan_session_id WHERE plan.story_id=? AND work.status<>'approved' ORDER BY work.position LIMIT 1`, storyID).Scan(&itemID, &state)
	if err != nil {
		return Session{}, err
	}
	if state != "failed" && state != "blocked" && state != "pending" {
		return Session{}, fmt.Errorf("work item is active or waiting for human review")
	}
	var currentKind SessionKind
	_ = tx.QueryRow("SELECT kind FROM sessions WHERE id=(SELECT current_session_id FROM work_items WHERE id=?)", itemID).Scan(&currentKind)
	kind := ImplementationSession
	if currentKind == SummarySession && summaryCurrent {
		kind = SummarySession
	}
	session, err := queueWorkItemSession(tx, itemID, kind)
	if err != nil {
		return Session{}, err
	}
	if err := tx.Commit(); err != nil {
		return Session{}, err
	}
	return session, nil
}
func (s *Store) assignedWorkItem(sessionID string) (WorkItem, error) {
	var id string
	if err := s.db.QueryRow("SELECT work_item_id FROM work_item_sessions WHERE session_id=?", sessionID).Scan(&id); err != nil {
		return WorkItem{}, err
	}
	return s.workItem(id)
}
func textList(values []string, maxCount int) ([]string, error) {
	if len(values) > maxCount {
		return nil, fmt.Errorf("too many entries")
	}
	result := []string{}
	for _, value := range values {
		value, err := progressSummary(value)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}
func compact(value string) string { return strings.TrimSpace(value) }
func progressSummary(text string) (string, error) {
	text = compact(text)
	if text == "" || !utf8.ValidString(text) || utf8.RuneCountInString(text) > 8000 {
		return "", fmt.Errorf("text must be 1–8000 characters")
	}
	return text, nil
}
