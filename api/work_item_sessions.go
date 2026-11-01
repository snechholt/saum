package main

import (
	"database/sql"
	"fmt"
)

func (s *Store) cancelQueuedSession(id string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.Exec("UPDATE sessions SET status='cancelled',finished_at=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=? AND status='queued'", id)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		if err := finishWorkItemSession(tx, id, SessionCancelled, "Stopped by operator"); err != nil {
			return 0, err
		}
	}
	return count, tx.Commit()
}
func (s *Store) sessionThread(session Session) (string, error) {
	if session.Kind == PlanningSession {
		return "", nil
	}
	var thread string
	var err error
	if session.Kind == ImplementationSession {
		err = s.db.QueryRow("SELECT thread_id FROM goal_workflows WHERE plan_session_id=?", session.ParentSessionID).Scan(&thread)
	} else {
		err = s.db.QueryRow("SELECT COALESCE((SELECT thread_id FROM factory_agents WHERE factory_id=story.factory_id),'') FROM user_stories AS story WHERE id=?", session.StoryID).Scan(&thread)
	}
	return thread, err
}
func (s *Store) saveSessionThread(session Session, thread string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("UPDATE sessions SET thread_id=? WHERE id=?", thread, session.ID); err != nil {
		return err
	}
	if session.Kind == ImplementationSession {
		if _, err := tx.Exec("UPDATE goal_workflows SET thread_id=? WHERE plan_session_id=?", thread, session.ParentSessionID); err != nil {
			return err
		}
	}
	if session.Kind == SummarySession {
		if _, err := tx.Exec("INSERT INTO factory_agents(factory_id,thread_id) SELECT factory_id,? FROM user_stories WHERE id=? ON CONFLICT(factory_id) DO UPDATE SET thread_id=excluded.thread_id", thread, session.StoryID); err != nil {
			return err
		}
	}
	return tx.Commit()
}
func finishWorkItemSession(tx *sql.Tx, id string, status SessionStatus, message string) error {
	var itemID, state string
	var kind SessionKind
	err := tx.QueryRow("SELECT work.id,work.status,run.kind FROM work_item_sessions AS assigned JOIN work_items AS work ON work.id=assigned.work_item_id JOIN sessions AS run ON run.id=assigned.session_id WHERE run.id=? AND work.current_session_id=run.id", id).Scan(&itemID, &state, &kind)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	if status == SessionCompleted {
		if kind == ImplementationSession && state != "review_requested" {
			return fmt.Errorf("work item did not request review")
		}
		if kind == SummarySession {
			var count int
			if err := tx.QueryRow("SELECT count(*) FROM work_item_reviews WHERE summary_session_id=?", id).Scan(&count); err != nil {
				return err
			}
			if count != 1 {
				return fmt.Errorf("factory agent did not submit an executive summary")
			}
			_, err = tx.Exec("UPDATE work_items SET status='awaiting_review' WHERE id=?", itemID)
			return err
		}
		return nil
	}
	if state == "blocked" {
		return nil
	}
	_, err = tx.Exec("UPDATE work_items SET status='failed',blocker=? WHERE id=?", message, itemID)
	return err
}
func (s *Store) workItemTurnReady(session Session) error {
	item, err := s.assignedWorkItem(session.ID)
	if err != nil {
		return err
	}
	if session.Kind == ImplementationSession {
		if item.Status == "blocked" {
			return fmt.Errorf("blocked: %s", item.Blocker)
		}
		if item.Status != "review_requested" {
			return fmt.Errorf("goal agent must submit each acceptance criterion and call saum.work_item.request_review()")
		}
		return nil
	}
	for _, review := range item.Reviews {
		if review.SummarySessionID == session.ID {
			return nil
		}
	}
	return fmt.Errorf("factory agent did not submit an executive summary")
}
