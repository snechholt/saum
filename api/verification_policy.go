package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type queryReader interface {
	Query(string, ...any) (*sql.Rows, error)
}

func requiredCommands(reader queryReader, planID string) ([][]string, error) {
	rows, err := reader.Query(`SELECT DISTINCT checks.command FROM work_item_reviews AS review
 JOIN work_items AS work ON work.id=review.work_item_id
 JOIN acceptance_submissions AS receipt ON receipt.session_id=review.goal_session_id AND receipt.work_item_id=work.id
 JOIN json_each(receipt.check_ids) AS reference JOIN verification_checks AS checks ON checks.id=reference.value
 WHERE work.plan_session_id=? AND review.decision='approved' ORDER BY checks.command`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	commands := [][]string{}
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		var command []string
		if err := json.Unmarshal([]byte(data), &command); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, rows.Err()
}
func requiredChecksPassed(tx *sql.Tx, planID, sessionID string, tree ArtifactTree) error {
	commands, err := requiredCommands(tx, planID)
	if err != nil {
		return err
	}
	for _, command := range commands {
		var count int
		if err := tx.QueryRow("SELECT count(*) FROM verification_checks WHERE session_id=? AND artifact_tree=? AND status='passed' AND exit_code=0 AND command=?", sessionID, tree, commandKey(command)).Scan(&count); err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("rerun the previously approved regression check on this source: %v", command)
		}
	}
	return nil
}
func (s *Store) checkReport(sessionID string) (any, error) {
	item, err := s.assignedWorkItem(sessionID)
	if err != nil {
		return nil, err
	}
	checks, err := s.verificationChecks(sessionID)
	if err != nil {
		return nil, err
	}
	if len(checks) > 20 {
		checks = checks[len(checks)-20:]
	}
	for index := range checks {
		if len(checks[index].Output) > 16000 {
			checks[index].Output = checks[index].Output[len(checks[index].Output)-16000:]
		}
	}
	commands, err := requiredCommands(s.db, item.PlanSessionID)
	if err != nil {
		return nil, err
	}
	return struct {
		Checks           []VerificationCheck `json:"checks"`
		RequiredCommands [][]string          `json:"requiredCommands"`
	}{checks, commands}, nil
}
