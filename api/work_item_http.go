package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func (s *Store) replyWorkItemHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		problem(w, 415, "Expected application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 40000)
	var input struct {
		Reply string `json:"reply"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		problem(w, 400, "Invalid reply")
		return
	}
	session, err := s.replyToWorkItem(r.PathValue("item"), input.Reply)
	if err != nil {
		problem(w, 409, err.Error())
		return
	}
	jsonResponse(w, 202, session)
}
func (s *Store) reworkIntegrationHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		problem(w, 415, "Expected application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 40000)
	var input struct {
		Feedback string `json:"feedback"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		problem(w, 400, "Invalid rework request")
		return
	}
	next, err := s.reworkIntegration(r.PathValue("item"), input.Feedback)
	if err != nil {
		problem(w, 409, err.Error())
		return
	}
	jsonResponse(w, 202, next)
}
func (s *Store) verificationHTTP(w http.ResponseWriter, r *http.Request) {
	checks, err := s.verificationChecks(r.PathValue("session"))
	if err != nil {
		problem(w, 500, "Could not load verification checks")
		return
	}
	jsonResponse(w, 200, checks)
}
func (s *Store) workItemDecisionHTTP(approve bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
			problem(w, 415, "Expected application/json")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 40000)
		var input WorkItemDecision
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF || input.ReviewID == "" {
			problem(w, 400, "Expected an executive review ID and optional feedback")
			return
		}
		session, err := s.decideWorkItem(r.PathValue("item"), input, approve)
		if errors.Is(err, sql.ErrNoRows) {
			problem(w, 404, "Work item review not found")
			return
		}
		if err != nil {
			problem(w, 409, err.Error())
			return
		}
		jsonResponse(w, 200, struct {
			Session *Session `json:"session"`
		}{session})
	}
}
