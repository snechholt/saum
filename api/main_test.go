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
func TestFactoriesPersistAndStartWithoutStories(t *testing.T) {
	file := filepath.Join(t.TempDir(), "factory.sqlite")
	store, err := openStore(file)
	if err != nil {
		t.Fatal(err)
	}
	factories, err := store.listFactories()
	if err != nil || len(factories) != 0 {
		t.Fatalf("new store not empty: %v %v", factories, err)
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
	stories, err := reopened.listStories(factory.ID)
	if err != nil || len(stories) != 0 {
		t.Fatalf("new factory has stories: %v %v", stories, err)
	}
}
func TestStoriesBelongToFactories(t *testing.T) {
	store := testStore(t)
	first, err := store.createFactory("One")
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.createFactory("Two")
	if err != nil {
		t.Fatal(err)
	}
	spec, err := parseStorySpec("Manage tasks", "People can keep track of what they need to do.")
	if err != nil {
		t.Fatal(err)
	}
	story, err := store.createStory(first.ID, spec)
	if err != nil {
		t.Fatal(err)
	}
	stories, err := store.listStories(first.ID)
	if err != nil || len(stories) != 1 || stories[0] != story {
		t.Fatalf("missing story: %v %v", stories, err)
	}
	stories, err = store.listStories(second.ID)
	if err != nil || len(stories) != 0 {
		t.Fatalf("story leaked across factories: %v %v", stories, err)
	}
	if _, err = store.db.Exec("INSERT INTO user_stories (id, factory_id, title) VALUES ('invalid', 'missing', 'Invalid')"); err == nil {
		t.Fatal("accepted orphan story")
	}
}
func TestImportPreservesEarlierEntriesOnce(t *testing.T) {
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
		stories, err := store.listStories("saum")
		if err != nil || len(stories) != 2 {
			t.Fatalf("import lost or duplicated entries: %v %v", stories, err)
		}
		_ = store.db.Close()
	}
	bytes, err := os.ReadFile(filepath.Join(dir, "projects.json"))
	if err != nil || string(bytes) != legacy {
		t.Fatal("original JSON changed")
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
		{"GET", "/api/factories/missing/stories", "", "", 404},
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
