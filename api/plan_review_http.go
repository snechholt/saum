package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func reviewProblem(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Plan or comment not found")
		return
	}
	if errors.Is(err, ErrReviewConflict) {
		problem(w, 409, err.Error())
		return
	}
	problem(w, 400, err.Error())
}
func (s *Store) postPlanComment(w http.ResponseWriter, r *http.Request) {
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		problem(w, 415, "Expected application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64000)
	var input CommentInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		problem(w, 400, "Invalid comment")
		return
	}
	comment, err := s.addPlanComment(r.PathValue("session"), input)
	if err != nil {
		reviewProblem(w, err)
		return
	}
	jsonResponse(w, 201, comment)
}
func (s *Store) deletePlanComment(w http.ResponseWriter, r *http.Request) {
	if err := s.removePlanComment(r.PathValue("session"), r.PathValue("comment")); err != nil {
		reviewProblem(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *Store) postPlanApproval(w http.ResponseWriter, r *http.Request) {
	if err := s.approvePlan(r.PathValue("session")); err != nil {
		reviewProblem(w, err)
		return
	}
	session, err := s.latestImplementation(r.PathValue("session"))
	if err != nil {
		problem(w, 500, "Could not load implementation session")
		return
	}
	jsonResponse(w, 200, session)
}
func (s *Store) postPlanFeedback(w http.ResponseWriter, r *http.Request) {
	session, err := s.sendPlanFeedback(r.PathValue("session"))
	if err != nil {
		reviewProblem(w, err)
		return
	}
	jsonResponse(w, 202, session)
}
