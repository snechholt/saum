package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFactoryCRUD(t *testing.T) {
	store := testStore(t)
	factory, err := store.createFactory("Original")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := store.updateFactory(factory.ID, " Renamed ")
	if err != nil || updated.ID != factory.ID || updated.Name != "Renamed" {
		t.Fatalf("update failed: %+v %v", updated, err)
	}
	if _, err = store.updateFactory(factory.ID, ""); err == nil {
		t.Fatal("accepted empty name")
	}
	if err = store.deleteFactory(factory.ID); err != nil {
		t.Fatal(err)
	}
	factories, err := store.listFactories()
	if err != nil || len(factories) != 0 {
		t.Fatalf("delete failed: %+v %v", factories, err)
	}
	if _, err = store.updateFactory(factory.ID, "Missing"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected missing factory, got %v", err)
	}
	if err = store.deleteFactory(factory.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected missing factory, got %v", err)
	}
}
func TestCannotDeleteFactoryWithProjects(t *testing.T) {
	store := testStore(t)
	factory, err := store.createFactory("Saum")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.db.Exec("INSERT INTO projects (id, factory_id, name) VALUES ('todo', ?, 'ToDo')", factory.ID); err != nil {
		t.Fatal(err)
	}
	if err = store.deleteFactory(factory.ID); !errors.Is(err, ErrFactoryNotEmpty) {
		t.Fatalf("expected nonempty factory error, got %v", err)
	}
	projects, err := store.listProjects(factory.ID)
	if err != nil || len(projects) != 1 {
		t.Fatalf("project lost: %+v %v", projects, err)
	}
}
func TestFactoryCRUDHTTP(t *testing.T) {
	store := testStore(t)
	empty, err := store.createFactory("Empty")
	if err != nil {
		t.Fatal(err)
	}
	nonempty, err := store.createFactory("Nonempty")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.db.Exec("INSERT INTO projects (id, factory_id, name) VALUES ('todo', ?, 'ToDo')", nonempty.ID); err != nil {
		t.Fatal(err)
	}
	handler := store.handler(http.NotFoundHandler())
	for _, tc := range []struct {
		method, id, body string
		want             int
	}{
		{"PUT", empty.ID, `{"name":"Renamed"}`, 200},
		{"PUT", empty.ID, `{"name":" "}`, 400},
		{"PUT", "missing", `{"name":"Missing"}`, 404},
		{"DELETE", nonempty.ID, "", 409},
		{"DELETE", empty.ID, "", 204},
		{"DELETE", empty.ID, "", 404},
	} {
		request := httptest.NewRequest(tc.method, "http://127.0.0.1:8080/api/factories/"+tc.id, strings.NewReader(tc.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != tc.want {
			t.Errorf("%s %s: got %d want %d: %s", tc.method, tc.id, response.Code, tc.want, response.Body.String())
		}
	}
}
