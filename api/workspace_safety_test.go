package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkspacePinCannotBeRedirectedAfterValidation(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	owned, err := s.ownedWorkspace(t.Context(), factoryID, plan.Worktree)
	if err != nil {
		t.Fatal(err)
	}
	pin, err := owned.pin()
	if err != nil {
		t.Fatal(err)
	}
	defer pin.Close()
	moved := plan.Worktree + ".moved"
	if err := os.Rename(plan.Worktree, moved); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, plan.Worktree); err != nil {
		t.Fatal(err)
	}
	command, err := directoryCommand(t.Context(), pin, "/bin/pwd")
	if err != nil {
		t.Fatal(err)
	}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatal(string(output), err)
	}
	expected, _ := filepath.EvalSymlinks(moved)
	if strings.TrimSpace(string(output)) != expected {
		t.Fatal("child escaped pinned directory", string(output))
	}
	if _, err := s.ownedWorkspace(t.Context(), factoryID, plan.Worktree); err == nil {
		t.Fatal("redirected workspace accepted")
	}
}
func TestRepositorySymlinkIsRejectedBeforeGitInitialization(t *testing.T) {
	s := testStore(t)
	factory, err := s.createFactory("Symlink probe")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(factory.WorkspaceRoot, 0700); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, factory.RepositoryPath); err != nil {
		t.Fatal(err)
	}
	planner := Planner{Store: s}
	if _, err := planner.worktree(t.Context(), factory.ID, "session"); err == nil {
		t.Fatal("initialized redirected repository")
	}
	if _, err := os.Stat(filepath.Join(outside, ".git")); !os.IsNotExist(err) {
		t.Fatal("outside repository modified", err)
	}
}
