package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func acceptedSibling(t *testing.T, local string) (*Store, OwnedWorkspace) {
	t.Helper()
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	p := Planner{Store: s}
	path, err := p.worktree(context.Background(), factoryID, "sibling-story")
	if err != nil {
		t.Fatal(err)
	}
	owned, err := s.ownedWorkspace(context.Background(), factoryID, path)
	if err != nil {
		t.Fatal(err)
	}
	if local != "" {
		if err := os.WriteFile(filepath.Join(path, local), []byte("local story work\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, s)
	if err := os.WriteFile(filepath.Join(plan.Worktree, "README.md"), []byte("accepted shared application\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, item := prepareSummary(t, s, goal)
	if _, err := s.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID}, true); err != nil {
		t.Fatal(err)
	}
	if err := s.processIntegrations(); err != nil {
		t.Fatal(err)
	}
	result, err := s.goalIntegration(plan.ID)
	if err != nil || result.Status != "completed" {
		t.Fatal(result, err)
	}
	return s, owned
}
func TestStoriesReceiveAcceptedMainAndPreserveUnapprovedWork(t *testing.T) {
	s, owned := acceptedSibling(t, "login.txt")
	note, err := s.syncStoryWorkspace(context.Background(), owned)
	if err != nil {
		t.Fatal(err)
	}
	for file, want := range map[string]string{"README.md": "accepted shared application\n", "login.txt": "local story work\n"} {
		data, err := os.ReadFile(filepath.Join(owned.Worktree, file))
		if err != nil || string(data) != want {
			t.Fatal(file, string(data), err)
		}
	}
	if !strings.Contains(note, "checkpoint") {
		t.Fatal(note)
	}
	if _, err := gitOutput(context.Background(), owned.Repository, "show", "HEAD:login.txt"); err == nil {
		t.Fatal("unapproved story work reached main")
	}
	if _, err := gitOutput(context.Background(), owned.Worktree, "merge-base", "--is-ancestor", "main", "HEAD"); err != nil {
		t.Fatal(err)
	}
	before, _ := gitOutput(context.Background(), owned.Worktree, "rev-parse", "HEAD")
	if _, err := s.syncStoryWorkspace(context.Background(), owned); err != nil {
		t.Fatal(err)
	}
	after, _ := gitOutput(context.Background(), owned.Worktree, "rev-parse", "HEAD")
	if string(before) != string(after) {
		t.Fatal("retry duplicated synchronization")
	}
}
func TestStorySynchronizationLeavesConflictsForAgentAndBlocksReview(t *testing.T) {
	s, owned := acceptedSibling(t, "README.md")
	note, err := s.syncStoryWorkspace(context.Background(), owned)
	if err != nil || !strings.Contains(note, "conflicts") {
		t.Fatal(note, err)
	}
	data, _ := os.ReadFile(filepath.Join(owned.Worktree, "README.md"))
	if !strings.Contains(string(data), "local story work") || !strings.Contains(string(data), "accepted shared application") {
		t.Fatal("merge lost source", string(data))
	}
	if _, err := captureTree(context.Background(), owned.Worktree, owned.Repository); err == nil {
		t.Fatal("unresolved merge accepted as verification artifact")
	}
	if _, err := s.syncStoryWorkspace(context.Background(), owned); err != nil {
		t.Fatal("retry cannot resume conflict resolution", err)
	}
}
func TestCleanStoryFastForwardsToAcceptedMain(t *testing.T) {
	s, owned := acceptedSibling(t, "")
	if _, err := s.syncStoryWorkspace(context.Background(), owned); err != nil {
		t.Fatal(err)
	}
	main, _ := gitOutput(context.Background(), owned.Repository, "rev-parse", "HEAD")
	head, _ := gitOutput(context.Background(), owned.Worktree, "rev-parse", "HEAD")
	if string(main) != string(head) {
		t.Fatal("clean story did not reuse factory main")
	}
}
