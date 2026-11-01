package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func successfulCheck(t *testing.T, s *Store, session Session, command []string) VerificationCheck {
	t.Helper()
	if session.ThreadID == "" {
		thread, err := s.sessionThread(session)
		if err != nil {
			t.Fatal(err)
		}
		if thread == "" {
			thread = "goal-thread"
		}
		if err := s.saveSessionThread(session, thread); err != nil {
			t.Fatal(err)
		}
	}
	check, err := s.requestVerification(session.ID, CheckRequest{Command: command})
	if err != nil {
		t.Fatal(err)
	}
	s.executeVerification(context.Background(), check)
	checks, err := s.verificationChecks(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range checks {
		if result.ID == check.ID {
			if result.Status != "passed" {
				t.Fatalf("check failed: %+v", result)
			}
			return result
		}
	}
	t.Fatal("check result missing")
	return VerificationCheck{}
}
func TestVerificationRecordsRealExitAndIsolatedSource(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	if err := os.WriteFile(filepath.Join(plan.Worktree, "value"), []byte("42"), 0600); err != nil {
		t.Fatal(err)
	}
	check := successfulCheck(t, s, goal, []string{"/bin/sh", "-c", "test \"$(cat value)\" = 42"})
	if check.ExitCode == nil || *check.ExitCode != 0 || check.ArtifactTree == "" {
		t.Fatal(check)
	}
	if _, err := os.Stat(filepath.Join(plan.Worktree, ".git")); err != nil {
		t.Fatal(err)
	}
	index := 0
	if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Actual check", CheckIDs: []string{check.ID}}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plan.Worktree, "value"), []byte("99"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Stale check", CheckIDs: []string{check.ID}}); err == nil {
		t.Fatal("old check accepted for new source")
	}
}
