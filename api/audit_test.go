package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditPagesAllSessionsWithoutMixingRolesOrFactories(t *testing.T) {
	store := testStore(t)
	plan, _ := submittedPlan(t, store)
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	factoryID, _, _, _ := store.storyContext(plan.StoryID)
	for index := range 110 {
		if _, err := store.db.Exec("INSERT INTO sessions(id,story_id,status,created_at) VALUES (?,?,'failed','2026-11-01T00:00:00Z')", fmt.Sprintf("history-%d", index), plan.StoryID); err != nil {
			t.Fatal(err)
		}
	}
	page, err := store.audit(0, "", "")
	if err != nil || len(page) != 100 {
		t.Fatal(len(page), err)
	}
	older, err := store.audit(page[len(page)-1].Cursor, "", "")
	if err != nil || len(older) != 12 {
		t.Fatal(len(older), err)
	}
	seen := map[string]bool{}
	for _, entry := range append(page, older...) {
		if seen[entry.Session.ID] {
			t.Fatal("duplicate page entry")
		}
		seen[entry.Session.ID] = true
		if entry.Factory.ID != factoryID || entry.Factory.WorkspaceRoot == "" {
			t.Fatal(entry)
		}
	}
	goals, err := store.audit(0, factoryID, ImplementationSession)
	if err != nil || len(goals) != 1 || goals[0].WorkItemTitle == "" {
		t.Fatal(goals, err)
	}
	missing, err := store.audit(0, "unknown-factory", "")
	if err != nil || len(missing) != 0 {
		t.Fatal(missing, err)
	}
	handler := store.handler(http.NotFoundHandler())
	for _, query := range []string{"before=-1", "before=invalid", "role=unknown"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", "http://127.0.0.1:8080/api/audit?"+query, nil))
		if response.Code != 400 {
			t.Fatal(query, response.Code)
		}
	}
}
