package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStoryDetailRetainsSessionHistory(t *testing.T) {
	store := testStore(t)
	factory, err := store.createFactory("Factory")
	if err != nil {
		t.Fatal(err)
	}
	spec, err := parseStorySpec("Sign in", "An existing user can access their private information.")
	if err != nil {
		t.Fatal(err)
	}
	story, err := store.createStory(factory.ID, spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.enqueueStories(); err != nil {
		t.Fatal(err)
	}
	first, err := store.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(first.ID, SessionFailed, "Provider unavailable"); err != nil {
		t.Fatal(err)
	}
	retry, err := store.retryStory(story.ID)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	store.handler(http.NotFoundHandler()).ServeHTTP(response, httptest.NewRequest("GET", "http://127.0.0.1:8080/api/stories/"+story.ID, nil))
	if response.Code != 200 {
		t.Fatal(response.Code, response.Body.String())
	}
	var result struct {
		Story    UserStory
		Factory  Factory
		Sessions []Session
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Factory.ID != factory.ID || result.Story.ID != story.ID || result.Story.SessionID != retry.ID || len(result.Sessions) != 2 || result.Sessions[1].ID != first.ID {
		t.Fatalf("invalid detail: %+v", result)
	}
	for _, path := range []string{"/api/stories/missing", "/api/sessions/missing", "/api/sessions/missing/events"} {
		response := httptest.NewRecorder()
		store.handler(http.NotFoundHandler()).ServeHTTP(response, httptest.NewRequest("GET", "http://127.0.0.1:8080"+path, nil))
		if response.Code != 404 {
			t.Fatalf("%s got %d", path, response.Code)
		}
	}
}
