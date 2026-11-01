package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type VerificationCheck struct {
	ID           string       `json:"id"`
	SessionID    string       `json:"sessionId"`
	WorkItemID   string       `json:"workItemId"`
	ArtifactTree ArtifactTree `json:"artifactTree"`
	Command      []string     `json:"command"`
	Status       string       `json:"status"`
	Output       string       `json:"output"`
	ExitCode     *int         `json:"exitCode"`
	Error        string       `json:"error"`
}
type CheckRequest struct {
	Command []string `json:"command"`
}

func parseCheckCommand(command []string) ([]string, error) {
	if len(command) == 0 || len(command) > 64 {
		return nil, fmt.Errorf("command must contain 1–64 arguments")
	}
	total := 0
	for _, argument := range command {
		total += len(argument)
		if strings.ContainsRune(argument, 0) {
			return nil, fmt.Errorf("command contains a null byte")
		}
	}
	if strings.TrimSpace(command[0]) == "" || strings.HasPrefix(command[0], "-") || total > 16000 {
		return nil, fmt.Errorf("invalid verification command")
	}
	return append([]string{}, command...), nil
}
func (s *Store) requestVerification(sessionID string, input CheckRequest) (VerificationCheck, error) {
	command, err := parseCheckCommand(input.Command)
	if err != nil {
		return VerificationCheck{}, err
	}
	data, _ := json.Marshal(command)
	tx, err := s.db.Begin()
	if err != nil {
		return VerificationCheck{}, err
	}
	defer tx.Rollback()
	item, err := assignedItem(tx, sessionID, ImplementationSession)
	if err != nil {
		return VerificationCheck{}, err
	}
	if item.Status != "working" {
		return VerificationCheck{}, fmt.Errorf("checks require a working item")
	}
	var pending int
	if err := tx.QueryRow("SELECT count(*) FROM verification_checks WHERE session_id=? AND status IN ('queued','running')", sessionID).Scan(&pending); err != nil {
		return VerificationCheck{}, err
	}
	var totalChecks int
	if err := tx.QueryRow("SELECT count(*) FROM verification_checks WHERE session_id=?", sessionID).Scan(&totalChecks); err != nil {
		return VerificationCheck{}, err
	}
	if totalChecks >= 100 {
		return VerificationCheck{}, fmt.Errorf("verification limit reached for this turn")
	}
	if pending >= 8 {
		return VerificationCheck{}, fmt.Errorf("too many pending checks")
	}
	check := VerificationCheck{ID: rand.Text(), SessionID: sessionID, WorkItemID: item.ID, Command: command, Status: "queued"}
	if _, err := tx.Exec("INSERT INTO verification_checks(id,session_id,work_item_id,artifact_tree,command,status) VALUES (?,?,?,'',?,'queued')", check.ID, sessionID, item.ID, string(data)); err != nil {
		return check, err
	}
	if err := toolJournal(tx, sessionID, "saum.work_item.run_check", input, check); err != nil {
		return check, err
	}
	return check, tx.Commit()
}
func (s *Store) verificationChecks(sessionID string) ([]VerificationCheck, error) {
	rows, err := s.db.Query("SELECT id,session_id,work_item_id,artifact_tree,command,status,output,exit_code,error FROM verification_checks WHERE session_id=? ORDER BY rowid", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	checks := []VerificationCheck{}
	for rows.Next() {
		var check VerificationCheck
		var command string
		if err := rows.Scan(&check.ID, &check.SessionID, &check.WorkItemID, &check.ArtifactTree, &command, &check.Status, &check.Output, &check.ExitCode, &check.Error); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(command), &check.Command); err != nil {
			return nil, err
		}
		checks = append(checks, check)
	}
	return checks, rows.Err()
}
func (s *Store) runVerificationChecks(ctx context.Context, session Session) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	defer s.db.Exec("UPDATE verification_checks SET status='cancelled',error='Agent turn ended before verification finished',finished_at=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE session_id=? AND status IN ('queued','running')", session.ID)
	for {
		if ctx.Err() != nil {
			return
		}
		checks, err := s.verificationChecks(session.ID)
		if err != nil {
			return
		}
		for _, check := range checks {
			if check.Status == "queued" {
				s.executeVerification(ctx, check)
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

type boundedOutput struct {
	mu   sync.Mutex
	data []byte
}

func (w *boundedOutput) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	size := len(data)
	if remaining := (512 << 10) - len(w.data); remaining > 0 {
		if len(data) > remaining {
			data = data[:remaining]
		}
		w.data = append(w.data, data...)
	}
	return size, nil
}
func (w *boundedOutput) String() string { w.mu.Lock(); defer w.mu.Unlock(); return string(w.data) }

func (s *Store) executeVerification(parent context.Context, check VerificationCheck) {
	ctx, cancel := context.WithTimeout(parent, 2*time.Minute)
	defer cancel()
	result, err := s.db.Exec("UPDATE verification_checks SET status='running' WHERE id=? AND status='queued' AND EXISTS(SELECT 1 FROM sessions WHERE id=session_id AND status='running')", check.ID)
	if err != nil {
		return
	}
	count, _ := result.RowsAffected()
	if count != 1 {
		return
	}
	finish := func(status string, output string, exit *int, problem error) {
		message := ""
		if problem != nil {
			message = problem.Error()
		}
		_, _ = s.db.Exec("UPDATE verification_checks SET status=?,output=?,exit_code=?,error=?,finished_at=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=?", status, output, exit, message, check.ID)
	}
	owned, tree, err := s.sessionArtifact(ctx, check.SessionID)
	if err != nil {
		finish("failed", "", nil, err)
		return
	}
	if _, err := s.db.Exec("UPDATE verification_checks SET artifact_tree=? WHERE id=?", tree, check.ID); err != nil {
		finish("failed", "", nil, err)
		return
	}
	area := filepath.Join(owned.Root, ".saum-checks")
	if err := os.MkdirAll(area, 0700); err != nil {
		finish("failed", "", nil, err)
		return
	}
	// OpenRoot prevents an ignored/symlinked path from redirecting the verifier.
	factory, err := os.OpenRoot(owned.Root)
	if err != nil {
		finish("failed", "", nil, err)
		return
	}
	defer factory.Close()
	if err := factory.Mkdir(filepath.Join(".saum-checks", check.ID), 0700); err != nil {
		finish("failed", "", nil, err)
		return
	}
	base := filepath.Join(area, check.ID)
	defer factory.RemoveAll(filepath.Join(".saum-checks", check.ID))
	source := filepath.Join(base, "source")
	if err := materializeTree(ctx, owned.Repository, tree, source); err != nil {
		finish("failed", "", nil, err)
		return
	}
	pinned, err := pinDirectory(source)
	if err != nil {
		finish("failed", "", nil, err)
		return
	}
	defer pinned.Close()
	command, err := directoryCommand(ctx, pinned, check.Command[0], check.Command[1:]...)
	if err != nil {
		finish("failed", "", nil, err)
		return
	}
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "SAUM_SESSION_TOKEN=") {
			command.Env = append(command.Env, entry)
		}
	}
	output := &boundedOutput{}
	command.Stdout = output
	command.Stderr = output
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.Cancel = func() error { return syscall.Kill(-command.Process.Pid, syscall.SIGKILL) }
	command.WaitDelay = time.Second
	if err := command.Start(); err != nil {
		finish("failed", output.String(), nil, err)
		return
	}
	identity := processIdentity(command.Process.Pid)
	for attempt := 0; attempt < 20 && isDirectoryTrampoline(identity); attempt++ {
		time.Sleep(5 * time.Millisecond)
		identity = processIdentity(command.Process.Pid)
	}
	if _, err := s.db.Exec("UPDATE verification_checks SET process_id=?,process_identity=? WHERE id=?", command.Process.Pid, identity, check.ID); err != nil {
		_ = command.Cancel()
		_ = command.Wait()
		finish("failed", output.String(), nil, err)
		return
	}
	waitErr := command.Wait()
	exit := command.ProcessState.ExitCode()
	if ctx.Err() != nil {
		finish("cancelled", output.String(), &exit, ctx.Err())
		return
	}
	if waitErr != nil {
		finish("failed", output.String(), &exit, waitErr)
		return
	}
	after, err := captureTree(ctx, source, owned.Repository)
	if err != nil {
		finish("failed", output.String(), &exit, err)
		return
	}
	if after != tree {
		finish("failed", output.String(), &exit, fmt.Errorf("verification changed source files; run checks without modifying the reviewed source"))
		return
	}
	finish("passed", output.String(), &exit, nil)
}
func verifiedCheckIDs(tx *sql.Tx, sessionID string, tree ArtifactTree, ids []string) error {
	if len(ids) == 0 || len(ids) > 20 {
		return fmt.Errorf("provide 1–20 recorded check IDs")
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			return fmt.Errorf("duplicate verification check ID")
		}
		seen[id] = true
		var status string
		var observed ArtifactTree
		var exit sql.NullInt64
		if err := tx.QueryRow("SELECT status,artifact_tree,exit_code FROM verification_checks WHERE id=? AND session_id=?", id, sessionID).Scan(&status, &observed, &exit); err != nil {
			return err
		}
		if status != "passed" || !exit.Valid || exit.Int64 != 0 || observed != tree {
			return fmt.Errorf("check %s did not pass against the current source snapshot", id)
		}
	}
	return nil
}
func (s *Store) recoverVerificationChecks() error {
	rows, err := s.db.Query("SELECT process_id,process_identity FROM verification_checks WHERE status='running'")
	if err != nil {
		return err
	}
	type process struct {
		pid      int
		identity string
	}
	processes := []process{}
	for rows.Next() {
		var p process
		if err := rows.Scan(&p.pid, &p.identity); err != nil {
			_ = rows.Close()
			return err
		}
		processes = append(processes, p)
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return err
	}
	for _, p := range processes {
		if p.pid > 0 && p.identity != "" && processIdentity(p.pid) == p.identity {
			_ = syscall.Kill(-p.pid, syscall.SIGKILL)
		}
	}
	_, err = s.db.Exec("UPDATE verification_checks SET status='cancelled',error='Server stopped during verification' WHERE status IN ('queued','running')")
	return err
}
func commandKey(command []string) string { data, _ := json.Marshal(command); return string(data) }
