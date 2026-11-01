package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Factory struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WorkspaceRoot  string `json:"workspaceRoot"`
	RepositoryPath string `json:"repositoryPath"`
}

// A UserStory describes a desired outcome within a Factory, not an implementation plan.
type UserStory struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Goal           string      `json:"goal"` // Empty only for entries imported before goals existed.
	PlanningStatus string      `json:"planningStatus"`
	SessionID      string      `json:"sessionId"`
	SessionKind    SessionKind `json:"sessionKind"`
	GoalStatus     string      `json:"goalStatus"`
}

var ErrFactoryNotEmpty = errors.New("Factory still contains user stories")

type Store struct {
	db      *sql.DB
	planner *Planner
}

func openStore(file string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	// Keep foreign-key settings and serialized writes on one connection.
	db.SetMaxOpenConns(1)
	if _, err = db.Exec("PRAGMA foreign_keys = ON; PRAGMA busy_timeout = 5000;"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err = migrate(db, filepath.Join(filepath.Dir(file), "projects.json")); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrateSettings(db, filepath.Dir(file)); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}
func (s *Store) listFactories() ([]Factory, error) {
	rows, err := s.db.Query("SELECT id,name,workspace_root,repository_path FROM factories ORDER BY rowid")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	factories := []Factory{}
	for rows.Next() {
		var factory Factory
		if err := rows.Scan(&factory.ID, &factory.Name, &factory.WorkspaceRoot, &factory.RepositoryPath); err != nil {
			return nil, err
		}
		factories = append(factories, factory)
	}
	return factories, rows.Err()
}
func (s *Store) createFactory(rawName string) (Factory, error) {
	name, err := parseName(rawName)
	if err != nil {
		return Factory{}, err
	}
	settings, err := s.settings()
	if err != nil {
		return Factory{}, err
	}
	factory := Factory{ID: rand.Text(), Name: name}
	factory.WorkspaceRoot = filepath.Join(settings.FactoriesRoot, factory.ID)
	factory.RepositoryPath = filepath.Join(factory.WorkspaceRoot, "main")
	_, err = s.db.Exec("INSERT INTO factories (id,name,workspace_root,repository_path) VALUES (?,?,?,?)", factory.ID, factory.Name, factory.WorkspaceRoot, factory.RepositoryPath)
	if err != nil {
		return Factory{}, err
	}
	return factory, nil
}
func (s *Store) listStories(factoryID string) ([]UserStory, error) {
	var exists int
	if err := s.db.QueryRow("SELECT 1 FROM factories WHERE id = ?", factoryID).Scan(&exists); err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`SELECT id, title, goal,
 CASE WHEN length(trim(goal))=0 THEN 'needs_goal' ELSE COALESCE((SELECT status FROM sessions WHERE story_id=user_stories.id ORDER BY rowid DESC LIMIT 1),'queued') END,
 COALESCE((SELECT id FROM sessions WHERE story_id=user_stories.id ORDER BY rowid DESC LIMIT 1),''),
 COALESCE((SELECT kind FROM sessions WHERE story_id=user_stories.id ORDER BY rowid DESC LIMIT 1),'planning'),
 CASE WHEN NOT EXISTS(SELECT 1 FROM goal_workflows JOIN sessions AS plan ON plan.id=plan_session_id WHERE plan.story_id=user_stories.id) THEN 'planning'
 WHEN EXISTS(SELECT 1 FROM work_items JOIN sessions AS plan ON plan.id=work_items.plan_session_id WHERE plan.story_id=user_stories.id AND work_items.status='blocked') THEN 'blocked'
 WHEN EXISTS(SELECT 1 FROM work_items JOIN sessions AS plan ON plan.id=work_items.plan_session_id WHERE plan.story_id=user_stories.id AND work_items.status='failed') THEN 'failed'
 WHEN EXISTS(SELECT 1 FROM work_items JOIN sessions AS plan ON plan.id=work_items.plan_session_id WHERE plan.story_id=user_stories.id AND work_items.status='awaiting_review') THEN 'awaiting_review'
 WHEN EXISTS(SELECT 1 FROM work_items JOIN sessions AS plan ON plan.id=work_items.plan_session_id WHERE plan.story_id=user_stories.id AND work_items.status<>'approved') THEN 'working'
 ELSE 'completed' END
 FROM user_stories WHERE factory_id = ? ORDER BY rowid`, factoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stories := []UserStory{}
	for rows.Next() {
		var story UserStory
		if err := rows.Scan(&story.ID, &story.Title, &story.Goal, &story.PlanningStatus, &story.SessionID, &story.SessionKind, &story.GoalStatus); err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}
	return stories, rows.Err()
}
