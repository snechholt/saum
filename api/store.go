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
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Projects have no behavior yet; they belong to a Factory.
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var ErrFactoryNotEmpty = errors.New("Factory still contains projects")

type Store struct{ db *sql.DB }

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
	return &Store{db: db}, nil
}
func (s *Store) listFactories() ([]Factory, error) {
	rows, err := s.db.Query("SELECT id, name FROM factories ORDER BY rowid")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	factories := []Factory{}
	for rows.Next() {
		var factory Factory
		if err := rows.Scan(&factory.ID, &factory.Name); err != nil {
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
	factory := Factory{ID: rand.Text(), Name: name}
	_, err = s.db.Exec("INSERT INTO factories (id, name) VALUES (?, ?)", factory.ID, factory.Name)
	if err != nil {
		return Factory{}, err
	}
	return factory, nil
}
func (s *Store) listProjects(factoryID string) ([]Project, error) {
	var exists int
	if err := s.db.QueryRow("SELECT 1 FROM factories WHERE id = ?", factoryID).Scan(&exists); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id, name FROM projects WHERE factory_id = ? ORDER BY rowid", factoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	projects := []Project{}
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}
