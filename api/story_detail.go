package main

import (
	"database/sql"
	"errors"
	"net/http"
)

func (s *Store) storySessions(storyID string) ([]Session, error) {
	rows, err := s.db.Query("SELECT "+sessionColumns+" FROM sessions WHERE story_id=? ORDER BY rowid DESC", storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sessions := []Session{}
	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}
func (s *Store) getStory(w http.ResponseWriter, r *http.Request) {
	factoryID, _, _, err := s.storyContext(r.PathValue("story"))
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "User story not found")
		return
	}
	if err != nil {
		problem(w, 500, "Could not load story")
		return
	}
	stories, err := s.listStories(factoryID)
	if err != nil {
		problem(w, 500, "Could not load story")
		return
	}
	var story UserStory
	for _, item := range stories {
		if item.ID == r.PathValue("story") {
			story = item
			break
		}
	}
	if story.ID == "" {
		problem(w, 404, "User story not found")
		return
	}
	var factory Factory
	if err := s.db.QueryRow("SELECT id,name,workspace_root,repository_path FROM factories WHERE id=?", factoryID).Scan(&factory.ID, &factory.Name, &factory.WorkspaceRoot, &factory.RepositoryPath); err != nil {
		problem(w, 500, "Could not load factory")
		return
	}
	sessions, err := s.storySessions(story.ID)
	if err != nil {
		problem(w, 500, "Could not load sessions")
		return
	}
	jsonResponse(w, 200, struct {
		Story    UserStory `json:"story"`
		Factory  Factory   `json:"factory"`
		Sessions []Session `json:"sessions"`
	}{story, factory, sessions})
}
