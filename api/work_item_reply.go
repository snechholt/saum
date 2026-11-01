package main

import "fmt"

func (s *Store) replyToWorkItem(itemID, reply string) (Session, error) {
	reply, err := progressSummary(reply)
	if err != nil {
		return Session{}, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return Session{}, err
	}
	defer tx.Rollback()
	var state string
	if err := tx.QueryRow("SELECT status FROM work_items WHERE id=?", itemID).Scan(&state); err != nil {
		return Session{}, err
	}
	if state != "blocked" {
		return Session{}, fmt.Errorf("work item is not blocked")
	}
	if _, err := tx.Exec("INSERT INTO work_item_messages(work_item_id,body) VALUES (?,?)", itemID, reply); err != nil {
		return Session{}, err
	}
	session, err := queueWorkItemSession(tx, itemID, ImplementationSession)
	if err != nil {
		return Session{}, err
	}
	if err := tx.Commit(); err != nil {
		return Session{}, err
	}
	return session, nil
}
