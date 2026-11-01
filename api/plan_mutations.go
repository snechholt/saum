package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

func mutablePlan(tx *sql.Tx, id string) error {
	var status SessionStatus
	var kind SessionKind
	var sealed bool
	if err := tx.QueryRow("SELECT status,kind,plan_submitted FROM sessions WHERE id=?", id).Scan(&status, &kind, &sealed); err != nil {
		return err
	}
	if kind != PlanningSession || status != SessionRunning || sealed {
		return ErrPlanClosed
	}
	return nil
}
func toolJournal(tx *sql.Tx, sessionID, name string, input, result any) error {
	data, err := json.Marshal(map[string]any{"tool": name, "input": input, "result": result})
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,'tool',?)", sessionID, data)
	return err
}
func (s *Store) updateItem(sessionID string, raw ItemSubmission) (ActionItem, error) {
	input, err := parseItem(raw)
	if err != nil {
		return ActionItem{}, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return ActionItem{}, err
	}
	defer tx.Rollback()
	if err := mutablePlan(tx, sessionID); err != nil {
		return ActionItem{}, err
	}
	criteria, err := json.Marshal(input.AcceptanceCriteria)
	if err != nil {
		return ActionItem{}, err
	}
	item := ActionItem{SessionID: sessionID, Key: input.Key, Title: input.Title, Description: input.Description, AcceptanceCriteria: input.AcceptanceCriteria}
	err = tx.QueryRow("UPDATE action_items SET title=?,description=?,acceptance_criteria=? WHERE session_id=? AND item_key=? RETURNING id", item.Title, item.Description, string(criteria), sessionID, item.Key).Scan(&item.ID)
	if err != nil {
		return ActionItem{}, err
	}
	if err := toolJournal(tx, sessionID, "saum.update_item", input, item); err != nil {
		return ActionItem{}, err
	}
	if err := tx.Commit(); err != nil {
		return ActionItem{}, err
	}
	return item, nil
}
func (s *Store) removeItem(sessionID, key string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := mutablePlan(tx, sessionID); err != nil {
		return err
	}
	result, err := tx.Exec("DELETE FROM action_items WHERE session_id=? AND item_key=?", sessionID, key)
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
	if err := toolJournal(tx, sessionID, "saum.remove_item", map[string]string{"key": key}, map[string]bool{"removed": true}); err != nil {
		return err
	}
	return tx.Commit()
}

type FeedbackResolution struct {
	CommentID  string `json:"commentId"`
	Resolution string `json:"resolution"`
}

func (s *Store) resolveFeedback(sessionID string, input FeedbackResolution) error {
	input.Resolution = strings.TrimSpace(input.Resolution)
	if !utf8.ValidString(input.Resolution) || input.Resolution == "" || utf8.RuneCountInString(input.Resolution) > 8000 {
		return fmt.Errorf("resolution must be 1–8000 characters")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := mutablePlan(tx, sessionID); err != nil {
		return err
	}
	var count int
	err = tx.QueryRow("SELECT count(*) FROM plan_comments WHERE id=? AND batch_id=(SELECT feedback_batch_id FROM sessions WHERE id=?)", input.CommentID, sessionID).Scan(&count)
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("feedback comment does not belong to this revision")
	}
	_, err = tx.Exec("INSERT INTO feedback_resolutions (session_id,comment_id,resolution) VALUES (?,?,?) ON CONFLICT(session_id,comment_id) DO UPDATE SET resolution=excluded.resolution", sessionID, input.CommentID, input.Resolution)
	if err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.resolve_feedback", input, map[string]bool{"resolved": true}); err != nil {
		return err
	}
	return tx.Commit()
}
