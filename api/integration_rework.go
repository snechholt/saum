package main

import (
	"database/sql"
	"fmt"
)

// A failed integration does not erase a human decision. Rework creates a new
// turn and later a new executive review; the approved review remains immutable.
func (s *Store) reworkIntegration(itemID, feedback string) (Session, error) {
	feedback, err := progressSummary(feedback)
	if err != nil {
		return Session{}, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return Session{}, err
	}
	defer tx.Rollback()
	var planID, state string
	var position int
	if err := tx.QueryRow("SELECT plan_session_id,status,position FROM work_items WHERE id=?", itemID).Scan(&planID, &state, &position); err != nil {
		return Session{}, err
	}
	var following, pending int
	if err := tx.QueryRow("SELECT count(*) FROM work_items WHERE plan_session_id=? AND position>?", planID, position).Scan(&following); err != nil {
		return Session{}, err
	}
	if err := tx.QueryRow("SELECT count(*) FROM work_items WHERE plan_session_id=? AND status<>'approved'", planID).Scan(&pending); err != nil {
		return Session{}, err
	}
	var integration string
	err = tx.QueryRow("SELECT status FROM integrations WHERE plan_session_id=? ORDER BY rowid DESC LIMIT 1", planID).Scan(&integration)
	if err != nil && err != sql.ErrNoRows {
		return Session{}, err
	}
	if state != "approved" || following != 0 || pending != 0 || integration == "queued" || integration == "running" || integration == "completed" {
		return Session{}, fmt.Errorf("only the final item of an unintegrated goal can be reworked")
	}
	if _, err := tx.Exec("INSERT INTO work_item_messages(work_item_id,body) VALUES (?,?)", itemID, feedback); err != nil {
		return Session{}, err
	}
	next, err := queueWorkItemSession(tx, itemID, ImplementationSession)
	if err != nil {
		return Session{}, err
	}
	return next, tx.Commit()
}
