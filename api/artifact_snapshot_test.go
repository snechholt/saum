package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotIncludesRawUntrackedBytesModesAndSafeLinks(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	_, repo, _ := s.factoryWorkspace(factoryID)
	name := "space quote\"\n雪.txt"
	data := []byte("raw\r\nsource\x00bytes")
	if err := os.WriteFile(filepath.Join(plan.Worktree, name), data, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(name, filepath.Join(plan.Worktree, "link")); err != nil {
		t.Fatal(err)
	}
	tree, err := captureTree(t.Context(), plan.Worktree, repo)
	if err != nil {
		t.Fatal(err)
	}
	copy := filepath.Join(t.TempDir(), "copy")
	if err := materializeTree(t.Context(), repo, tree, copy); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(copy, name))
	if err != nil || string(got) != string(data) {
		t.Fatal("snapshot transformed source bytes", err)
	}
	after, err := captureTree(t.Context(), copy, repo)
	if err != nil || after != tree {
		t.Fatal("snapshot did not round-trip", tree, after, err)
	}
	external := t.TempDir()
	if err := os.Symlink(external, filepath.Join(plan.Worktree, "escape")); err != nil {
		t.Fatal(err)
	}
	if _, err := captureTree(t.Context(), plan.Worktree, repo); err == nil {
		t.Fatal("snapshot followed a source symlink outside its root")
	}
}
