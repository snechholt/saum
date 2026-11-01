package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type SessionKind string

const (
	PlanningSession       SessionKind = "planning"
	ImplementationSession SessionKind = "implementation"
)

type ImplementationItem struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Summary string `json:"summary"`
}
type Implementation struct {
	PlanSessionID string               `json:"planSessionId"`
	Submitted     bool                 `json:"submitted"`
	Summary       string               `json:"summary"`
	Items         []ImplementationItem `json:"items"`
}

func migrateImplementations(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`ALTER TABLE sessions ADD COLUMN kind TEXT NOT NULL DEFAULT 'planning' CHECK(kind IN ('planning','implementation'));
 CREATE TABLE implementations (
 session_id TEXT PRIMARY KEY NOT NULL REFERENCES sessions(id),plan_session_id TEXT NOT NULL REFERENCES sessions(id),
 submitted INTEGER NOT NULL DEFAULT 0,summary TEXT NOT NULL DEFAULT ''
 );
 CREATE INDEX builds_by_plan ON implementations(plan_session_id);
 CREATE TABLE implementation_items (
 session_id TEXT NOT NULL REFERENCES implementations(session_id),plan_item_id TEXT NOT NULL REFERENCES action_items(id),item_key TEXT NOT NULL,
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','running','completed','blocked')),
 summary TEXT NOT NULL DEFAULT '',PRIMARY KEY(session_id,item_key)
 );
 CREATE UNIQUE INDEX one_running_item ON implementation_items(session_id) WHERE status='running';
 PRAGMA user_version=6;`)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func queueImplementation(tx *sql.Tx, planID string) (Session, error) {
	var storyID string
	var review ReviewStatus
	var status SessionStatus
	var sealed bool
	var kind SessionKind
	if err := tx.QueryRow("SELECT story_id,review_status,status,plan_submitted,kind FROM sessions WHERE id=?", planID).Scan(&storyID, &review, &status, &sealed, &kind); err != nil {
		return Session{}, err
	}
	if kind != PlanningSession || status != SessionCompleted || !sealed || review != ReviewApproved {
		return Session{}, fmt.Errorf("implementation requires an approved, submitted plan")
	}
	session := Session{ID: rand.Text(), StoryID: storyID, Status: SessionQueued, Kind: ImplementationSession, ParentSessionID: planID, ReviewStatus: ReviewPending, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if _, err := tx.Exec("INSERT INTO sessions (id,story_id,status,kind,parent_session_id,created_at) VALUES (?,?,'queued','implementation',?,?)", session.ID, storyID, planID, session.CreatedAt); err != nil {
		return Session{}, err
	}
	if _, err := tx.Exec("INSERT INTO implementations (session_id,plan_session_id) VALUES (?,?)", session.ID, planID); err != nil {
		return Session{}, err
	}
	result, err := tx.Exec("INSERT INTO implementation_items (session_id,plan_item_id,item_key) SELECT ?,id,item_key FROM action_items WHERE session_id=? ORDER BY rowid", session.ID, planID)
	if err != nil {
		return Session{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Session{}, err
	}
	if count == 0 {
		return Session{}, fmt.Errorf("approved plan is empty")
	}
	event, _ := json.Marshal(map[string]any{"type": "implementation_queued", "approvedPlanId": planID})
	if _, err := tx.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,'workflow',?)", session.ID, event); err != nil {
		return Session{}, err
	}
	return session, nil
}
func (s *Store) queueApprovedPlans() error {
	// Catch up approvals made before the build transition existed, without repeating any attempts.
	rows, err := s.db.Query(`SELECT id FROM sessions AS plan WHERE kind='planning' AND review_status='approved'
 AND NOT EXISTS(SELECT 1 FROM implementations WHERE plan_session_id=plan.id)
 AND NOT EXISTS(SELECT 1 FROM sessions WHERE story_id=plan.story_id AND status IN ('queued','preparing','running'))
 AND id=(SELECT id FROM sessions WHERE story_id=plan.story_id AND kind='planning' ORDER BY rowid DESC LIMIT 1)`)
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
	err = errors.Join(rows.Err(), rows.Close())
	if err != nil {
		return err
	}
	for _, id := range ids {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		var count int
		err = tx.QueryRow("SELECT count(*) FROM implementations WHERE plan_session_id=?", id).Scan(&count)
		if err == nil && count == 0 {
			_, err = queueImplementation(tx, id)
		}
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) latestImplementation(planID string) (Session, error) {
	return scanSession(s.db.QueryRow("SELECT "+sessionColumns+" FROM sessions WHERE id=(SELECT session_id FROM implementations WHERE plan_session_id=? ORDER BY rowid DESC LIMIT 1)", planID))
}
func (s *Store) retryImplementation(storyID string) (Session, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Session{}, err
	}
	defer tx.Rollback()
	var planID string
	if err := tx.QueryRow("SELECT id FROM sessions WHERE story_id=? AND kind='planning' AND review_status='approved' ORDER BY rowid DESC LIMIT 1", storyID).Scan(&planID); err != nil {
		return Session{}, err
	}
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM sessions WHERE story_id=? AND status IN ('queued','preparing','running')", storyID).Scan(&count); err != nil {
		return Session{}, err
	}
	if count > 0 {
		return Session{}, fmt.Errorf("user story already has an active session")
	}
	session, err := queueImplementation(tx, planID)
	if err != nil {
		return Session{}, err
	}
	if err := tx.Commit(); err != nil {
		return Session{}, err
	}
	return session, nil
}
func (s *Store) implementation(id string) (*Implementation, error) {
	var build Implementation
	err := s.db.QueryRow("SELECT plan_session_id,submitted,summary FROM implementations WHERE session_id=?", id).Scan(&build.PlanSessionID, &build.Submitted, &build.Summary)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT work.item_key,plan.title,work.status,work.summary FROM implementation_items AS work JOIN action_items AS plan ON plan.id=work.plan_item_id WHERE work.session_id=? ORDER BY work.rowid", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	build.Items = []ImplementationItem{}
	for rows.Next() {
		var item ImplementationItem
		if err := rows.Scan(&item.Key, &item.Title, &item.Status, &item.Summary); err != nil {
			return nil, err
		}
		build.Items = append(build.Items, item)
	}
	return &build, rows.Err()
}
func implementationOpen(tx *sql.Tx, id string) error {
	var kind SessionKind
	var status SessionStatus
	var submitted bool
	var approved ReviewStatus
	err := tx.QueryRow(`SELECT run.kind,run.status,build.submitted,plan.review_status FROM sessions AS run
 JOIN implementations AS build ON build.session_id=run.id JOIN sessions AS plan ON plan.id=build.plan_session_id WHERE run.id=?`, id).Scan(&kind, &status, &submitted, &approved)
	if err != nil {
		return err
	}
	if kind != ImplementationSession || status != SessionRunning || submitted || approved != ReviewApproved {
		return fmt.Errorf("implementation is not accepting progress updates")
	}
	return nil
}
func (s *Store) startImplementationItem(sessionID, key string) (ImplementationItem, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return ImplementationItem{}, err
	}
	defer tx.Rollback()
	if err := implementationOpen(tx, sessionID); err != nil {
		return ImplementationItem{}, err
	}
	var status string
	if err := tx.QueryRow("SELECT status FROM implementation_items WHERE session_id=? AND item_key=?", sessionID, key).Scan(&status); err != nil {
		return ImplementationItem{}, err
	}
	if status == "completed" {
		return ImplementationItem{}, fmt.Errorf("item is already completed")
	}
	var other int
	if err := tx.QueryRow("SELECT count(*) FROM implementation_items WHERE session_id=? AND status='running' AND item_key<>?", sessionID, key).Scan(&other); err != nil {
		return ImplementationItem{}, err
	}
	if other > 0 {
		return ImplementationItem{}, fmt.Errorf("finish the current item before starting another")
	}
	var item ImplementationItem
	item.Key = key
	item.Status = "running"
	if err := tx.QueryRow("UPDATE implementation_items SET status='running',summary='' WHERE session_id=? AND item_key=? RETURNING (SELECT title FROM action_items WHERE id=plan_item_id)", sessionID, key).Scan(&item.Title); err != nil {
		return ImplementationItem{}, err
	}
	if err := toolJournal(tx, sessionID, "saum.start_item", map[string]string{"key": key}, item); err != nil {
		return ImplementationItem{}, err
	}
	if err := tx.Commit(); err != nil {
		return ImplementationItem{}, err
	}
	return item, nil
}

type ItemProgress struct {
	Key     string `json:"key"`
	Summary string `json:"summary"`
}

func progressSummary(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" || !utf8.ValidString(text) || utf8.RuneCountInString(text) > 8000 {
		return "", fmt.Errorf("summary must be 1–8000 characters")
	}
	return text, nil
}
func (s *Store) reportImplementationItem(sessionID string, input ItemProgress, blocked bool) error {
	summary, err := progressSummary(input.Summary)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := implementationOpen(tx, sessionID); err != nil {
		return err
	}
	var state, previous string
	if err := tx.QueryRow("SELECT status,summary FROM implementation_items WHERE session_id=? AND item_key=?", sessionID, input.Key).Scan(&state, &previous); err != nil {
		return err
	}
	next, name := "completed", "saum.complete_item"
	if blocked {
		next, name = "blocked", "saum.block_item"
	}
	if state == next && previous == summary {
		return tx.Commit()
	}
	if state != "running" {
		return fmt.Errorf("start the item before reporting its result")
	}
	if _, err := tx.Exec("UPDATE implementation_items SET status=?,summary=? WHERE session_id=? AND item_key=?", next, summary, sessionID, input.Key); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, name, input, map[string]string{"status": next}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) completeImplementation(sessionID, rawSummary string) error {
	summary, err := progressSummary(rawSummary)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := implementationOpen(tx, sessionID); err != nil {
		return err
	}
	var remaining int
	if err := tx.QueryRow("SELECT count(*) FROM implementation_items WHERE session_id=? AND status<>'completed'", sessionID).Scan(&remaining); err != nil {
		return err
	}
	if remaining > 0 {
		return fmt.Errorf("%d plan items are not complete", remaining)
	}
	if _, err := tx.Exec("UPDATE implementations SET submitted=1,summary=? WHERE session_id=?", summary, sessionID); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.complete_implementation", map[string]string{"summary": summary}, map[string]bool{"submitted": true}); err != nil {
		return err
	}
	return tx.Commit()
}
