package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFactoryDefaultsAreSnapshotsNotRelocations(t *testing.T) {
	store := testStore(t)
	first, err := store.createFactory("First")
	if err != nil {
		t.Fatal(err)
	}
	originalRoot, originalRepo := first.WorkspaceRoot, first.RepositoryPath
	settings, err := store.settings()
	if err != nil || settings.Theme != "system" || !filepath.IsAbs(settings.FactoriesRoot) {
		t.Fatal(settings, err)
	}
	nextRoot := filepath.Join(t.TempDir(), "new-factories")
	if _, err := store.saveSettings(Settings{FactoriesRoot: nextRoot, Theme: "dark"}); err != nil {
		t.Fatal(err)
	}
	updated, err := store.updateFactory(first.ID, "Renamed")
	if err != nil || updated.WorkspaceRoot != originalRoot || updated.RepositoryPath != originalRepo {
		t.Fatal("settings or rename moved the factory", updated, err)
	}
	second, err := store.createFactory("Second")
	if err != nil || second.WorkspaceRoot != filepath.Join(nextRoot, second.ID) || second.RepositoryPath != filepath.Join(second.WorkspaceRoot, "main") {
		t.Fatal(second, err)
	}
	planner := Planner{Store: store}
	path, err := planner.worktree(t.Context(), second.ID, "goal-worktree")
	if err != nil {
		t.Fatal(err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(second.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(canonicalRoot, "goal-worktree") {
		t.Fatal("worktree not beside main", path)
	}
	if _, err := os.Stat(filepath.Join(second.RepositoryPath, ".git")); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(filepath.Join(path, ".git")); err != nil || info.IsDir() {
		t.Fatal("not a linked Git worktree", err)
	}
	if _, err := os.Stat(originalRoot); !os.IsNotExist(err) {
		t.Fatal("settings save touched the old workspace", err)
	}
}
func TestSettingsRejectInvalidBoundaries(t *testing.T) {
	store := testStore(t)
	valid := t.TempDir()
	file := filepath.Join(valid, "file")
	_ = os.WriteFile(file, []byte("not a directory"), 0600)
	for _, input := range []Settings{{FactoriesRoot: "../relative", Theme: "system"}, {FactoriesRoot: "/", Theme: "light"}, {FactoriesRoot: file, Theme: "dark"}, {FactoriesRoot: valid, Theme: "unknown"}} {
		if _, err := store.saveSettings(input); err == nil {
			t.Fatal("accepted invalid settings", input)
		}
	}
	handler := store.handler(http.NotFoundHandler())
	request := httptest.NewRequest("PUT", "http://127.0.0.1:8080/api/settings", strings.NewReader(`{"factoriesRoot":"/tmp/factories","theme":"light","moveExisting":true}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 400 {
		t.Fatal("unknown setting accepted", response.Code)
	}
}
