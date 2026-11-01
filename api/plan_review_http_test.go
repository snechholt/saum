package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestPlanReviewHTTP(t *testing.T) {
	store := testStore(t)
	source, item := submittedPlan(t, store)
	handler := store.handler(http.NotFoundHandler())
	input, _ := json.Marshal(titleComment(item))
	request := httptest.NewRequest("POST", "http://127.0.0.1:8080/api/sessions/"+source.ID+"/comments", bytes.NewReader(input))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 201 {
		t.Fatal(response.Code, response.Body.String())
	}
	var comment PlanComment
	if err := json.Unmarshal(response.Body.Bytes(), &comment); err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest("POST", "http://127.0.0.1:8080/api/sessions/"+source.ID+"/approve", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 409 {
		t.Fatal("approval did not reject pending comments", response.Code)
	}
	request = httptest.NewRequest("GET", "http://127.0.0.1:8080/api/sessions/"+source.ID, nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var detail struct{ Comments []PlanComment }
	if err := json.Unmarshal(response.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	if len(detail.Comments) != 1 || detail.Comments[0].ID != comment.ID {
		t.Fatal("saved comments not returned to browser")
	}
	request = httptest.NewRequest("DELETE", "http://127.0.0.1:8080/api/sessions/"+source.ID+"/comments/"+comment.ID, nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 204 {
		t.Fatal(response.Code, response.Body.String())
	}
	request = httptest.NewRequest("POST", "http://127.0.0.1:8080/api/sessions/"+source.ID+"/approve", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatal(response.Code, response.Body.String())
	}
	var approved Session
	if err := json.Unmarshal(response.Body.Bytes(), &approved); err != nil || approved.Kind != ImplementationSession || approved.ParentSessionID != source.ID {
		t.Fatal("approval not returned", err)
	}
}
func TestConcurrentFeedbackSendsQueueOnlyOneSession(t *testing.T) {
	store := testStore(t)
	source, item := submittedPlan(t, store)
	if _, err := store.addPlanComment(source.ID, titleComment(item)); err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	results := make(chan Session, 2)
	failures := make(chan error, 2)
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			session, err := store.sendPlanFeedback(source.ID)
			if err != nil {
				failures <- err
			} else {
				results <- session
			}
		}()
	}
	group.Wait()
	close(results)
	close(failures)
	for err := range failures {
		t.Fatal(err)
	}
	var id string
	count := 0
	for session := range results {
		if id != "" && session.ID != id {
			t.Fatal("feedback dispatched more than once")
		}
		id = session.ID
		count++
	}
	if count != 2 {
		t.Fatal("both deliveries should return the same session")
	}
}
