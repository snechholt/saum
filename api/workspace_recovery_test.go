package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func legacyApprovedFixture(t *testing.T) (*Store, Session) {
	t.Helper()
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	if err := os.WriteFile(filepath.Join(plan.Worktree, "shared.go"), []byte("package main\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, item := prepareSummary(t, s, goal)
	if _, err := s.db.Exec("UPDATE work_item_reviews SET decision='approved',artifact_tree='',approved_commit='' WHERE work_item_id=?", item.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec("UPDATE work_items SET status='approved' WHERE id=?", item.ID); err != nil {
		t.Fatal(err)
	}
	return s, plan
}
func TestLegacyRecoveryIntegratesWithoutForgingHistoricalEvidence(t *testing.T) {
	s, plan := legacyApprovedFixture(t)
	if err := s.recoverApprovedStory(context.Background(), plan.ID); err != nil {
		t.Fatal(err)
	}
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	_, repo, _ := s.factoryWorkspace(factoryID)
	if _, err := os.Stat(filepath.Join(repo, "shared.go")); err != nil {
		t.Fatal(err)
	}
	workflow, err := s.workflow(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	review := workflow.Items[0].Reviews[0]
	if review.ArtifactTree != "" || review.ApprovedCommit != "" || review.Decision != ReviewApproved {
		t.Fatal("legacy evidence rewritten", review)
	}
	before, _ := gitOutput(context.Background(), repo, "rev-parse", "HEAD")
	if err := s.recoverApprovedStory(context.Background(), plan.ID); err != nil {
		t.Fatal(err)
	}
	after, _ := gitOutput(context.Background(), repo, "rev-parse", "HEAD")
	if string(before) != string(after) {
		t.Fatal("recovery is not idempotent")
	}
}
func TestLegacyRecoveryRejectsUnapprovedWork(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.recoverApprovedStory(context.Background(), plan.ID); err == nil {
		t.Fatal("repair bypassed human approval")
	}
}

func TestRejectedBootstrapIsArchivedWithoutResumingOrApproving(t *testing.T) {
	s, accepted := legacyApprovedFixture(t)
	factoryID, _, _, _ := s.storyContext(accepted.StoryID)
	spec, err := parseStorySpec("Login", "Add login to the shared application")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.createStory(factoryID, spec); err != nil {
		t.Fatal(err)
	}
	if err := s.enqueueStories(); err != nil {
		t.Fatal(err)
	}
	plan, err := s.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	p := Planner{Store: s}
	path, err := p.worktree(context.Background(), factoryID, plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.startSession(plan.ID, path, "plan", "token"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.submitItem(plan.ID, ItemSubmission{Key: "login", Title: "Login", Description: "Add login", AcceptanceCriteria: []string{"Login works"}}); err != nil {
		t.Fatal(err)
	}
	if err := s.completePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.finishSession(plan.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	if err := os.WriteFile(filepath.Join(path, "app.py"), []byte("rejected independent app"), 0600); err != nil {
		t.Fatal(err)
	}
	_, item := prepareSummary(t, s, goal)
	if _, err := s.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID, Feedback: "Use the shared Go server"}, false); err != nil {
		t.Fatal(err)
	}
	resumed := startClaimed(t, s)
	if err := s.blockWorkItem(resumed.ID, "Shared server is missing"); err != nil {
		t.Fatal(err)
	}
	if err := s.finishSession(resumed.ID, SessionFailed, "Blocked"); err != nil {
		t.Fatal(err)
	}
	if err := s.recoverApprovedStory(context.Background(), accepted.ID); err != nil {
		t.Fatal(err)
	}
	backup, err := s.archiveRejectedBootstrap(context.Background(), plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(backup, "app.py"))
	if err != nil || string(data) != "rejected independent app" {
		t.Fatal("rejected code lost", err)
	}
	if _, err := os.Stat(filepath.Join(path, "shared.go")); err != nil {
		t.Fatal(err)
	}
	fresh, err := s.workItem(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fresh.Status != "blocked" || fresh.Reviews[0].Decision != ReviewChangesRequested {
		t.Fatal("repair bypassed approval", fresh)
	}
	if err := s.requireOfflineRepair(); err != nil {
		t.Fatal("repair queued work", err)
	}
}
