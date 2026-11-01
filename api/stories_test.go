package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestStorySpecRequiresAnOutcome(t *testing.T) {
	for _, input := range []struct{ title, goal string }{{"", "Useful outcome"}, {"Title", " "}, {"Title", strings.Repeat("x", 4001)}, {"Title", "invalid\x00goal"}} {
		if _, err := parseStorySpec(input.title, input.goal); err == nil {
			t.Errorf("accepted invalid story: %+v", input)
		}
	}
	spec, err := parseStorySpec("  Manage tasks ", "  People can capture tasks.\nThey can find outstanding work.  ")
	if err != nil || spec.title != "Manage tasks" || spec.goal != "People can capture tasks.\nThey can find outstanding work." {
		t.Fatalf("unexpected parsed story: %+v %v", spec, err)
	}
	if _, err := parseStorySpec("Title", strings.Repeat("🎯", 4000)); err != nil {
		t.Fatal("limits should count characters, not bytes", err)
	}
}
func TestCreateUserStoryHTTP(t *testing.T) {
	store := testStore(t)
	factory, err := store.createFactory("Saum")
	if err != nil {
		t.Fatal(err)
	}
	handler := store.handler(http.NotFoundHandler())
	for _, tc := range []struct {
		id, body string
		want     int
	}{
		{factory.ID, `{"title":"Manage tasks","goal":"People can capture and complete their daily tasks."}`, 201},
		{factory.ID, `{"title":"Missing goal"}`, 400},
		{factory.ID, `{"title":" ","goal":"Useful outcome"}`, 400},
		{factory.ID, `{"title":"Title","goal":"Useful outcome","framework":"React"}`, 400},
		{"missing", `{"title":"Title","goal":"Useful outcome"}`, 404},
	} {
		request := httptest.NewRequest("POST", "http://127.0.0.1:8080/api/factories/"+tc.id+"/stories", strings.NewReader(tc.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != tc.want {
			t.Errorf("got %d want %d: %s", response.Code, tc.want, response.Body.String())
		}
		if response.Code == 201 {
			var story UserStory
			if err := json.Unmarshal(response.Body.Bytes(), &story); err != nil {
				t.Fatal(err)
			}
			if story.ID == "" || story.Title != "Manage tasks" || story.Goal == "" {
				t.Fatalf("invalid story: %+v", story)
			}
		}
	}
	stories, err := store.listStories(factory.ID)
	if err != nil || len(stories) != 1 {
		t.Fatalf("invalid inputs created stories: %+v %v", stories, err)
	}
	if err := store.deleteFactory(factory.ID); err != ErrFactoryNotEmpty {
		t.Fatalf("story did not protect its factory: %v", err)
	}
}
func TestEarlierSchemasMigrateWithoutLosingEntries(t *testing.T) {
	for _, version := range []int{1, 2} {
		file := filepath.Join(t.TempDir(), "factory.sqlite")
		db, err := sql.Open("sqlite", file)
		if err != nil {
			t.Fatal(err)
		}
		schema := `CREATE TABLE factories (id TEXT PRIMARY KEY, name TEXT NOT NULL);
CREATE TABLE projects (id TEXT PRIMARY KEY, factory_id TEXT NOT NULL REFERENCES factories(id), name TEXT NOT NULL);
INSERT INTO factories VALUES ('factory', 'Saum'); INSERT INTO projects VALUES ('story', 'factory', 'Manage tasks'); PRAGMA user_version = 1;`
		if _, err = db.Exec(schema); err != nil {
			t.Fatal(err)
		}
		if version == 2 {
			if _, err = db.Exec("ALTER TABLE projects ADD COLUMN goal TEXT NOT NULL DEFAULT ''; UPDATE projects SET goal = 'People can track daily tasks'; PRAGMA user_version = 2;"); err != nil {
				t.Fatal(err)
			}
		}
		_ = db.Close()
		for range 2 {
			store, err := openStore(file)
			if err != nil {
				t.Fatal(err)
			}
			stories, err := store.listStories("factory")
			if err != nil || len(stories) != 1 || stories[0].ID != "story" || stories[0].Title != "Manage tasks" {
				t.Fatalf("lost entry from v%d: %+v %v", version, stories, err)
			}
			if version == 2 && stories[0].Goal != "People can track daily tasks" {
				t.Fatal("lost goal")
			}
			var oldTables int
			if err := store.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE name = 'projects'").Scan(&oldTables); err != nil || oldTables != 0 {
				t.Fatal("old project model remains", err)
			}
			_ = store.db.Close()
		}
	}
}
