package main

import (
	"crypto/rand"
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Factory struct{ db *sql.DB }

func openFactory(file string) (*Factory, error) {
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	// One connection serializes the small local factory's reads and writes.
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
  id TEXT PRIMARY KEY NOT NULL,
  name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 80)
 )`)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Factory{db: db}, nil
}
func (f *Factory) list() ([]Project, error) {
	rows, err := f.db.Query("SELECT id, name FROM projects ORDER BY rowid")
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
func (f *Factory) create(rawName string) (Project, error) {
	name, err := parseName(rawName)
	if err != nil {
		return Project{}, err
	}
	project := Project{ID: rand.Text(), Name: name}
	_, err = f.db.Exec("INSERT INTO projects (id, name) VALUES (?, ?)", project.ID, project.Name)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
