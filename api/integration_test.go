package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationRefusesDirtyMainAndPreservesApproval(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	if err := os.WriteFile(filepath.Join(plan.Worktree, "app.txt"), []byte("reviewed implementation"), 0600); err != nil {
		t.Fatal(err)
	}
	_, item := prepareSummary(t, s, goal)
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	_, repo, _ := s.factoryWorkspace(factoryID)
	if err := os.WriteFile(filepath.Join(repo, "operator.txt"), []byte("do not overwrite"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID}, true); err != nil {
		t.Fatal(err)
	}
	stories, _ := s.listStories(factoryID)
	if stories[0].GoalStatus != "integrating" {
		t.Fatal("goal advertised completion before integration", stories)
	}
	if err := s.processIntegrations(); err != nil {
		t.Fatal(err)
	}
	integration, err := s.goalIntegration(plan.ID)
	if err != nil || integration.Status != "failed" || !strings.Contains(integration.Error, "local changes") {
		t.Fatal(integration, err)
	}
	data, _ := os.ReadFile(filepath.Join(repo, "operator.txt"))
	if string(data) != "do not overwrite" {
		t.Fatal("operator data changed")
	}
	if _, err := os.Stat(filepath.Join(repo, "app.txt")); !os.IsNotExist(err) {
		t.Fatal("dirty main received code", err)
	}
	fresh, _ := s.workItem(item.ID)
	if fresh.Reviews[0].Decision != ReviewApproved || fresh.Reviews[0].ApprovedCommit == "" {
		t.Fatal("approval history was lost", fresh)
	}
	next, err := s.reworkIntegration(item.ID, "Rebase and reverify before integration.")
	if err != nil || next.Kind != ImplementationSession {
		t.Fatal(next, err)
	}
	fresh, _ = s.workItem(item.ID)
	if fresh.Reviews[0].Decision != ReviewApproved || fresh.Status != "queued" {
		t.Fatal("rework overwrote the approved review", fresh)
	}
}
func TestChecksCannotBeBorrowedAcrossTurnsOrSurviveSourceChanges(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	check := successfulCheck(t, s, goal, []string{"/usr/bin/true"})
	for index := range 2 {
		if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Recorded", CheckIDs: []string{check.ID}}); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(plan.Worktree, "changed.txt"), []byte("new source"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := s.requestWorkItemReview(goal.ID); err == nil {
		t.Fatal("stale receipts reached review")
	}
	updated := successfulCheck(t, s, goal, []string{"/usr/bin/true"})
	index := 0
	if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Reverified", CheckIDs: []string{updated.ID}}); err != nil {
		t.Fatal(err)
	}
	if err := s.requestWorkItemReview(goal.ID); err == nil {
		t.Fatal("later stale criterion was not invalidated by source version")
	}
	index = 1
	if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Reverified", CheckIDs: []string{updated.ID}}); err != nil {
		t.Fatal(err)
	}
	if err := s.requestWorkItemReview(goal.ID); err != nil {
		t.Fatal(err)
	}
}
func TestInternalSymlinkCannotTurnGoalIntoMain(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	_, repo, _ := s.factoryWorkspace(factoryID)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	run, err := s.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(plan.Worktree, plan.Worktree+".original"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(repo, plan.Worktree); err != nil {
		t.Fatal(err)
	}
	p := Planner{Store: s}
	if _, err := p.sessionWorkspace(t.Context(), run, factoryID); err == nil {
		t.Fatal("goal was granted the main checkout")
	}
}
