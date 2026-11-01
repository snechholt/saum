package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestFactoryStartsEmptyAndPersistsNamedProjects(t *testing.T) {
	file := filepath.Join(t.TempDir(), "projects.json")
	factory, err := openFactory(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(factory.projects) != 0 {
		t.Fatal("new factory should be empty")
	}
	project, err := factory.create("ToDo App")
	if err != nil {
		t.Fatal(err)
	}
	reopened, err := openFactory(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(reopened.projects) != 1 || reopened.projects[0] != project {
		t.Fatalf("project was not persisted: %+v", reopened.projects)
	}
}
func TestProjectNames(t *testing.T) {
	for _, name := range []string{"", "  ", "two\nlines", strings.Repeat("x", 81)} {
		if _, err := parseName(name); err == nil {
			t.Errorf("accepted %q", name)
		}
	}
	name, err := parseName("  ToDo App  ")
	if err != nil || name != "ToDo App" {
		t.Fatalf("unexpected name: %q %v", name, err)
	}
}
func TestProjectHTTPBoundary(t *testing.T) {
	factory, err := openFactory(filepath.Join(t.TempDir(), "projects.json"))
	if err != nil {
		t.Fatal(err)
	}
	handler := factory.handler(http.NotFoundHandler())
	for _, tc := range []struct {
		method, body, origin string
		want                 int
	}{
		{"GET", "", "", 200},
		{"POST", `{"name":"ToDo App"}`, "http://127.0.0.1:5173", 201},
		{"POST", `{"name":" "}`, "", 400},
		{"POST", `{"name":"test","issues":[]}`, "", 400},
		{"POST", `{"name":"test"} {}`, "", 400},
		{"POST", `{"name":"test"}`, "https://attacker.example", 403},
	} {
		request := httptest.NewRequest(tc.method, "http://127.0.0.1:8080/api/projects", strings.NewReader(tc.body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", tc.origin)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != tc.want {
			t.Errorf("got %d want %d: %s", recorder.Code, tc.want, recorder.Body.String())
		}
		if recorder.Code == 201 {
			var project Project
			if err := json.Unmarshal(recorder.Body.Bytes(), &project); err != nil {
				t.Fatal(err)
			}
			if project.ID == "" || project.Name != "ToDo App" {
				t.Fatalf("unexpected project: %+v", project)
			}
		}
	}
}
