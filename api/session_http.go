package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

func (s *Store) getSession(w http.ResponseWriter, r *http.Request) {
	session, err := s.session(r.PathValue("session"))
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Session not found")
		return
	}
	if err != nil {
		problem(w, 500, "Could not load session")
		return
	}
	build, err := s.implementation(session.ID)
	if err != nil {
		problem(w, 500, "Could not load implementation")
		return
	}
	planSession := session
	if build != nil {
		planSession, err = s.session(build.PlanSessionID)
		if err != nil {
			problem(w, 500, "Could not load approved plan")
			return
		}
	}
	items, err := s.items(planSession.ID)
	if err != nil {
		problem(w, 500, "Could not load items")
		return
	}
	comments, err := s.planComments(planSession.ID)
	if err != nil {
		problem(w, 500, "Could not load comments")
		return
	}
	feedback, err := s.receivedFeedback(session.FeedbackBatchID)
	if err != nil {
		problem(w, 500, "Could not load feedback")
		return
	}
	jsonResponse(w, 200, struct {
		Session          Session         `json:"session"`
		PlanSession      Session         `json:"planSession"`
		Implementation   *Implementation `json:"implementation"`
		Items            []ActionItem    `json:"items"`
		Comments         []PlanComment   `json:"comments"`
		ReceivedFeedback []PlanComment   `json:"receivedFeedback"`
	}{session, planSession, build, items, comments, feedback})
}
func (s *Store) getSessionEvents(w http.ResponseWriter, r *http.Request) {
	if _, err := s.session(r.PathValue("session")); errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Session not found")
		return
	} else if err != nil {
		problem(w, 500, "Could not load session")
		return
	}
	after := int64(0)
	if raw := r.URL.Query().Get("after"); raw != "" {
		var err error
		after, err = strconv.ParseInt(raw, 10, 64)
		if err != nil || after < 0 {
			problem(w, 400, "Invalid event cursor")
			return
		}
	}
	events, err := s.events(r.PathValue("session"), after)
	if err != nil {
		problem(w, 500, "Could not load session events")
		return
	}
	jsonResponse(w, 200, events)
}
func (s *Store) retryPlanning(w http.ResponseWriter, r *http.Request) {
	session, err := s.retryStory(r.PathValue("story"))
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "User story not found")
		return
	}
	if err != nil {
		problem(w, 409, err.Error())
		return
	}
	jsonResponse(w, 202, session)
}
func (s *Store) cancelSession(w http.ResponseWriter, r *http.Request) {
	session, err := s.session(r.PathValue("session"))
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Session not found")
		return
	}
	if err != nil {
		problem(w, 500, "Could not load session")
		return
	}
	if session.Status == SessionQueued {
		result, err := s.db.Exec("UPDATE sessions SET status='cancelled',finished_at=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=? AND status='queued'", session.ID)
		if err != nil {
			problem(w, 500, "Could not cancel session")
			return
		}
		count, _ := result.RowsAffected()
		if count == 0 && (s.planner == nil || !s.planner.Cancel(session.ID)) {
			problem(w, 409, "Session is being claimed; try again")
			return
		}
	} else if session.Status == SessionPreparing || session.Status == SessionRunning {
		if s.planner == nil || !s.planner.Cancel(session.ID) {
			problem(w, 409, "Session is being claimed; try again")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
