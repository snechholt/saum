package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	store, err := openStore(filepath.Join(t.TempDir(), "factory.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })
	return store
}
func TestFactoriesPersistAndStartWithoutProjects(t *testing.T) {
	file := filepath.Join(t.TempDir(), "factory.sqlite")
	store, err := openStore(file)
	if err != nil {
		t.Fatal(err)
	}
	factories, err := store.listFactories()
	if err != nil || len(factories) != 0 {
		t.Fatalf("new store is not empty: %v %v", factories, err)
	}
	factory, err := store.createFactory("  Saum  ")
	if err != nil {
		t.Fatal(err)
	}
	if factory.Name != "Saum" {
		t.Fatal("name not normalized")
	}
	if err := store.db.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := openStore(file)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.db.Close()
	factories, err = reopened.listFactories()
	if err != nil || len(factories) != 1 || factories[0] != factory {
		t.Fatalf("factory not persisted: %v %v", factories, err)
	}
	projects, err := reopened.listProjects(factory.ID)
	if err != nil || len(projects) != 0 {
		t.Fatalf("new factory has projects: %v %v", projects, err)
	}
}
func TestProjectsBelongToFactories(t *testing.T) {
	store := testStore(t)
	first, err := store.createFactory("One")
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.createFactory("Two")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.db.Exec("INSERT INTO projects (id, factory_id, name) VALUES ('todo', ?, 'ToDo')", first.ID); err != nil {
		t.Fatal(err)
	}
	projects, err := store.listProjects(first.ID)
	if err != nil || len(projects) != 1 || projects[0].Name != "ToDo" {
		t.Fatalf("missing project: %v %v", projects, err)
	}
	projects, err = store.listProjects(second.ID)
	if err != nil || len(projects) != 0 {
		t.Fatalf("project leaked across factories: %v %v", projects, err)
	}
	if _, err = store.db.Exec("INSERT INTO projects (id, factory_id, name) VALUES ('invalid', 'missing', 'Invalid')"); err == nil {
		t.Fatal("accepted orphan project")
	}
}
func TestImportPreservesPreviousProjectsOnce(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "factory.sqlite")
	db, err := sql.Open("sqlite", file)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE projects (id TEXT PRIMARY KEY, name TEXT NOT NULL); INSERT INTO projects VALUES ('existing', 'Existing')")
	if err != nil {
		t.Fatal(err)
	}
	_ = db.Close()
	legacy := `[{"id":"existing","name":"Existing"},{"id":"todo","name":"ToDo"}]`
	if err = os.WriteFile(filepath.Join(dir, "projects.json"), []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		store, err := openStore(file)
		if err != nil {
			t.Fatal(err)
		}
		factories, err := store.listFactories()
		if err != nil || len(factories) != 1 || factories[0].Name != "Saum" {
			t.Fatalf("missing imported factory: %v %v", factories, err)
		}
		projects, err := store.listProjects("saum")
		if err != nil || len(projects) != 2 {
			t.Fatalf("import lost or duplicated projects: %v %v", projects, err)
		}
		_ = store.db.Close()
	}
	bytes, err := os.ReadFile(filepath.Join(dir, "projects.json"))
	if err != nil || string(bytes) != legacy {
		t.Fatal("original JSON was changed")
	}
}
func TestNames(t *testing.T) {
	store := testStore(t)
	for _, name := range []string{"", "  ", "two\nlines", strings.Repeat("x", 81)} {
		if _, err := store.createFactory(name); err == nil {
			t.Errorf("accepted %q", name)
		}
	}
}
func TestFactoryHTTPBoundary(t *testing.T) {
	store := testStore(t)
	handler := store.handler(http.NotFoundHandler())
	for _, tc := range []struct {
		method, path, body, origin string
		want                       int
	}{
		{"GET", "/api/factories", "", "", 200},
		{"POST", "/api/factories", `{"name":"Saum"}`, "http://127.0.0.1:5173", 201},
		{"POST", "/api/factories", `{"name":"Other"}`, "http://localhost:5173", 201},
		{"POST", "/api/factories", `{"name":" "}`, "", 400},
		{"POST", "/api/factories", `{"name":"test","issues":[]}`, "", 400},
		{"POST", "/api/factories", `{"name":"test"} {}`, "", 400},
		{"POST", "/api/factories", `{"name":"test"}`, "https://attacker.example", 403},
		{"GET", "/api/factories/missing/projects", "", "", 404},
		{"POST", "/api/projects", `{"name":"test"}`, "", 404},
	} {
		request := httptest.NewRequest(tc.method, "http://127.0.0.1:8080"+tc.path, strings.NewReader(tc.body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", tc.origin)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != tc.want {
			t.Errorf("%s %s: got %d want %d: %s", tc.method, tc.path, response.Code, tc.want, response.Body.String())
		}
		if response.Code == 201 {
			var factory Factory
			if err := json.Unmarshal(response.Body.Bytes(), &factory); err != nil {
				t.Fatal(err)
			}
			if factory.ID == "" || factory.Name == "" {
				t.Fatalf("invalid factory: %+v", factory)
			}
		}
	}
}
