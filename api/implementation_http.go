package main

import (
	"database/sql"
	"errors"
	"net/http"
)

func (s *Store) retryImplementationHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := s.retryImplementation(r.PathValue("story"))
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "No approved plan found")
		return
	}
	if err != nil {
		problem(w, 409, err.Error())
		return
	}
	jsonResponse(w, 202, session)
}
