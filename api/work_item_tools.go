package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func assignedItem(tx *sql.Tx, sessionID string, kind SessionKind) (WorkItem, error) {
	var id string
	var status SessionStatus
	var actualKind SessionKind
	if err := tx.QueryRow("SELECT work_item_id,status,kind FROM work_item_sessions JOIN sessions ON sessions.id=session_id WHERE session_id=?", sessionID).Scan(&id, &status, &actualKind); err != nil {
		return WorkItem{}, err
	}
	if status != SessionRunning || actualKind != kind {
		return WorkItem{}, fmt.Errorf("function is not available in this session phase")
	}
	item, err := scanWorkItem(tx.QueryRow("SELECT "+workItemColumns+" FROM work_items AS work JOIN action_items AS item ON item.id=work.id WHERE work.id=?", id))
	if err != nil {
		return item, err
	}
	if item.CurrentSessionID != sessionID {
		return item, fmt.Errorf("session is no longer assigned to this work item")
	}
	return item, nil
}

type CriterionInput struct {
	CriterionIndex *int     `json:"criterionIndex"`
	Summary        string   `json:"summary"`
	Evidence       []string `json:"evidence"`
}

func (s *Store) submitAcceptanceCriterion(sessionID string, input CriterionInput) error {
	if input.CriterionIndex == nil {
		return fmt.Errorf("criterionIndex is required")
	}
	summary, err := progressSummary(input.Summary)
	if err != nil {
		return err
	}
	evidence, err := textList(input.Evidence, 20)
	if err != nil {
		return err
	}
	if len(evidence) == 0 {
		return fmt.Errorf("verification evidence is required")
	}
	data, _ := json.Marshal(evidence)
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	item, err := assignedItem(tx, sessionID, ImplementationSession)
	if err != nil {
		return err
	}
	if item.Status != "working" {
		return fmt.Errorf("item is not accepting criterion submissions")
	}
	index := *input.CriterionIndex
	if index < 0 || index >= len(item.AcceptanceCriteria) {
		return fmt.Errorf("criterionIndex is outside the assigned item's criteria")
	}
	var previousSummary, previousEvidence string
	err = tx.QueryRow("SELECT summary,evidence FROM acceptance_submissions WHERE session_id=? AND criterion_index=?", sessionID, index).Scan(&previousSummary, &previousEvidence)
	if err == nil {
		if previousSummary != summary || previousEvidence != string(data) {
			return fmt.Errorf("criterion was already submitted with different evidence")
		}
		return tx.Commit()
	}
	if err != sql.ErrNoRows {
		return err
	}
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM acceptance_submissions WHERE session_id=?", sessionID).Scan(&count); err != nil {
		return err
	}
	if index != count {
		return fmt.Errorf("submit criteria in order; next criterion is %d", count)
	}
	if _, err := tx.Exec("INSERT INTO acceptance_submissions (session_id,work_item_id,criterion_index,summary,evidence) VALUES (?,?,?,?,?)", sessionID, item.ID, index, summary, string(data)); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.submit_acceptance_criteria", input, map[string]bool{"submitted": true}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) requestWorkItemReview(sessionID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	item, err := assignedItem(tx, sessionID, ImplementationSession)
	if err != nil {
		return err
	}
	if item.Status == "review_requested" {
		return tx.Commit()
	}
	if item.Status != "working" {
		return fmt.Errorf("work item is not accepting review requests")
	}
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM acceptance_submissions WHERE session_id=?", sessionID).Scan(&count); err != nil {
		return err
	}
	if count != len(item.AcceptanceCriteria) || count == 0 {
		return fmt.Errorf("submit evidence for every acceptance criterion before requesting review")
	}
	if _, err := tx.Exec("UPDATE work_items SET status='review_requested' WHERE id=?", item.ID); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.request_review", map[string]any{}, map[string]bool{"requested": true}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) blockWorkItem(sessionID, reason string) error {
	reason, err := progressSummary(reason)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	item, err := assignedItem(tx, sessionID, ImplementationSession)
	if err != nil {
		return err
	}
	if item.Status != "working" && item.Status != "blocked" {
		return fmt.Errorf("cannot block an item after requesting review")
	}
	if _, err := tx.Exec("UPDATE work_items SET status='blocked',blocker=? WHERE id=?", reason, item.ID); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.report_blocker", map[string]string{"reason": reason}, map[string]bool{"blocked": true}); err != nil {
		return err
	}
	return tx.Commit()
}

type ExecutiveSummaryInput struct {
	Summary      string   `json:"summary"`
	Verification string   `json:"verification"`
	Risks        []string `json:"risks"`
}

func (s *Store) submitExecutiveSummary(sessionID string, input ExecutiveSummaryInput) error {
	summary, err := progressSummary(input.Summary)
	if err != nil {
		return err
	}
	verification, err := progressSummary(input.Verification)
	if err != nil {
		return err
	}
	risks, err := textList(input.Risks, 30)
	if err != nil {
		return err
	}
	risksJSON, _ := json.Marshal(risks)
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	item, err := assignedItem(tx, sessionID, SummarySession)
	if err != nil {
		return err
	}
	if item.Status != "summarizing" {
		return fmt.Errorf("item is not awaiting a factory-agent summary")
	}
	var previousSummary, previousVerification, previousRisks string
	err = tx.QueryRow("SELECT summary,verification,risks FROM work_item_reviews WHERE summary_session_id=?", sessionID).Scan(&previousSummary, &previousVerification, &previousRisks)
	if err == nil {
		if previousSummary != summary || previousVerification != verification || previousRisks != string(risksJSON) {
			return fmt.Errorf("executive summary was already submitted")
		}
		return tx.Commit()
	}
	if err != sql.ErrNoRows {
		return err
	}
	if _, err := tx.Exec("INSERT INTO work_item_reviews(id,work_item_id,goal_session_id,summary_session_id,summary,verification,risks) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?)", item.ID, item.LastGoalSessionID, sessionID, summary, verification, string(risksJSON)); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.submit_executive_summary", input, map[string]bool{"submitted": true}); err != nil {
		return err
	}
	return tx.Commit()
}
