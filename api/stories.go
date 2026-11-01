package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

type StorySpec struct{ title, goal string }

func parseStorySpec(rawTitle, rawGoal string) (StorySpec, error) {
	title, err := parseName(rawTitle)
	if err != nil {
		return StorySpec{}, fmt.Errorf("Title must contain 1–80 characters without control characters")
	}
	goal := strings.TrimSpace(rawGoal)
	if !utf8.ValidString(goal) || utf8.RuneCountInString(goal) < 1 || utf8.RuneCountInString(goal) > 4000 {
		return StorySpec{}, fmt.Errorf("Goal must contain 1–4000 characters")
	}
	for _, character := range goal {
		if unicode.IsControl(character) && character != '\n' && character != '\r' && character != '\t' {
			return StorySpec{}, fmt.Errorf("Goal contains unsupported control characters")
		}
	}
	return StorySpec{title: title, goal: goal}, nil
}
func (s *Store) createStory(factoryID string, spec StorySpec) (UserStory, error) {
	story := UserStory{PlanningStatus: "queued", SessionKind: PlanningSession}
	err := s.db.QueryRow(`INSERT INTO user_stories (id, factory_id, title, goal)
 SELECT ?, id, ?, ? FROM factories WHERE id = ?
 RETURNING id, title, goal`, rand.Text(), spec.title, spec.goal, factoryID).Scan(&story.ID, &story.Title, &story.Goal)
	return story, err
}
func (s *Store) postStory(w http.ResponseWriter, r *http.Request) {
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		problem(w, 415, "Expected application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20000)
	var input struct {
		Title string `json:"title"`
		Goal  string `json:"goal"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		problem(w, 400, "Invalid user story request")
		return
	}
	spec, err := parseStorySpec(input.Title, input.Goal)
	if err != nil {
		problem(w, 400, err.Error())
		return
	}
	story, err := s.createStory(r.PathValue("factory"), spec)
	if errors.Is(err, sql.ErrNoRows) {
		problem(w, 404, "Factory not found")
		return
	}
	if err != nil {
		problem(w, 500, "Could not save user story")
		return
	}
	jsonResponse(w, 201, story)
}
