package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func assignedItem(tx *sql.Tx, sessionID string, kind SessionKind) (WorkItem, error) {
	var id string
	var status SessionStatus
	var actualKind SessionKind
	var thread string
	if err := tx.QueryRow("SELECT work_item_id,status,kind,thread_id FROM work_item_sessions JOIN sessions ON sessions.id=session_id WHERE session_id=?", sessionID).Scan(&id, &status, &actualKind, &thread); err != nil {
		return WorkItem{}, err
	}
	if thread == "" {
		return WorkItem{}, fmt.Errorf("native thread identity is not established")
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
	CheckIDs       []string `json:"checkIds"`
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
	data, _ := json.Marshal(evidence)
	if len(input.CheckIDs) == 0 {
		return fmt.Errorf("recorded check IDs are required; use saum.work_item.run_check and checks")
	}
	_, tree, err := s.sessionArtifact(context.Background(), sessionID)
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
	if item.Status != "working" {
		return fmt.Errorf("item is not accepting criterion submissions")
	}
	index := *input.CriterionIndex
	if index < 0 || index >= len(item.AcceptanceCriteria) {
		return fmt.Errorf("criterionIndex is outside the assigned item")
	}
	if err := verifiedCheckIDs(tx, sessionID, tree, input.CheckIDs); err != nil {
		return err
	}
	checks, _ := json.Marshal(input.CheckIDs)
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM acceptance_submissions WHERE session_id=?", sessionID).Scan(&count); err != nil {
		return err
	}
	if index > count {
		return fmt.Errorf("submit criteria in order; next criterion is %d", count)
	}
	var previousSummary, previousEvidence, previousChecks string
	var previousTree ArtifactTree
	err = tx.QueryRow("SELECT summary,evidence,artifact_tree,check_ids FROM acceptance_submissions WHERE session_id=? AND criterion_index=?", sessionID, index).Scan(&previousSummary, &previousEvidence, &previousTree, &previousChecks)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && previousSummary == summary && previousEvidence == string(data) && previousTree == tree && previousChecks == string(checks) {
		return tx.Commit()
	}
	if _, err := tx.Exec("INSERT INTO acceptance_history(session_id,work_item_id,criterion_index,summary,evidence,artifact_tree,check_ids) VALUES (?,?,?,?,?,?,?)", sessionID, item.ID, index, summary, string(data), tree, string(checks)); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO acceptance_submissions(session_id,work_item_id,criterion_index,summary,evidence,artifact_tree,check_ids) VALUES (?,?,?,?,?,?,?) ON CONFLICT(session_id,criterion_index) DO UPDATE SET summary=excluded.summary,evidence=excluded.evidence,artifact_tree=excluded.artifact_tree,check_ids=excluded.check_ids", sessionID, item.ID, index, summary, string(data), tree, string(checks)); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.submit_acceptance_criteria", input, map[string]any{"submitted": true, "artifactTree": tree}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) requestWorkItemReview(sessionID string) error {
	owned, tree, err := s.sessionArtifact(context.Background(), sessionID)
	if err != nil {
		return err
	}
	session, err := s.session(sessionID)
	if err != nil {
		return err
	}
	if err := s.initializeGoalBase(context.Background(), session, owned); err != nil {
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
	if item.Status == "review_requested" {
		var sealed ArtifactTree
		if err := tx.QueryRow("SELECT review_tree FROM work_items WHERE id=?", item.ID).Scan(&sealed); err != nil {
			return err
		}
		if sealed != tree {
			return fmt.Errorf("source changed after requesting review")
		}
		return tx.Commit()
	}
	if item.Status != "working" {
		return fmt.Errorf("item is not accepting review requests")
	}
	var pending int
	if err := tx.QueryRow("SELECT count(*) FROM verification_checks WHERE session_id=? AND status IN ('queued','running')", sessionID).Scan(&pending); err != nil {
		return err
	}
	if pending > 0 {
		return fmt.Errorf("wait for pending checks before requesting review")
	}
	rows, err := tx.Query("SELECT criterion_index,artifact_tree,check_ids FROM acceptance_submissions WHERE session_id=? ORDER BY criterion_index", sessionID)
	if err != nil {
		return err
	}
	type receipt struct {
		index  int
		tree   ArtifactTree
		checks []string
	}
	receipts := []receipt{}
	for rows.Next() {
		var r receipt
		var ids string
		if err := rows.Scan(&r.index, &r.tree, &ids); err != nil {
			_ = rows.Close()
			return err
		}
		if err := json.Unmarshal([]byte(ids), &r.checks); err != nil {
			_ = rows.Close()
			return err
		}
		receipts = append(receipts, r)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if len(receipts) != len(item.AcceptanceCriteria) || len(receipts) == 0 {
		return fmt.Errorf("submit evidence for every acceptance criterion before requesting review")
	}
	for index, r := range receipts {
		if r.index != index || r.tree != tree {
			return fmt.Errorf("criterion %d needs reverification against the current source", index)
		}
		if err := verifiedCheckIDs(tx, sessionID, tree, r.checks); err != nil {
			return err
		}
	}
	if err := requiredChecksPassed(tx, item.PlanSessionID, sessionID, tree); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE work_items SET status='review_requested',review_tree=? WHERE id=?", tree, item.ID); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.request_review", map[string]any{}, map[string]any{"requested": true, "artifactTree": tree}); err != nil {
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
	_, tree, err := s.sessionArtifact(context.Background(), sessionID)
	if err != nil {
		return err
	}
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
	var sealed ArtifactTree
	if err := tx.QueryRow("SELECT review_tree FROM work_items WHERE id=?", item.ID).Scan(&sealed); err != nil {
		return err
	}
	if sealed == "" || sealed != tree {
		return fmt.Errorf("source changed since the review request; return to the goal agent for verification")
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
	if _, err := tx.Exec("INSERT INTO work_item_reviews(id,work_item_id,goal_session_id,summary_session_id,summary,verification,risks,artifact_tree) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?,?)", item.ID, item.LastGoalSessionID, sessionID, summary, verification, string(risksJSON), tree); err != nil {
		return err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.submit_executive_summary", input, map[string]bool{"submitted": true}); err != nil {
		return err
	}
	return tx.Commit()
}
