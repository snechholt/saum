package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// One-time cutover from the earlier unscoped project list. Original JSON is kept.
func migrate(db *sql.DB, legacyFile string) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version == 1 {
		return nil
	}
	if version != 0 {
		return fmt.Errorf("unsupported database version %d", version)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	projects, err := previousProjects(tx, legacyFile)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`CREATE TABLE factories (
 id TEXT PRIMARY KEY NOT NULL,
 name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 80)
 );
 CREATE TABLE projects_next (
 id TEXT PRIMARY KEY NOT NULL,
 factory_id TEXT NOT NULL REFERENCES factories(id),
 name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 80)
 );`)
	if err != nil {
		return err
	}
	if len(projects) > 0 {
		if _, err = tx.Exec("INSERT INTO factories (id, name) VALUES ('saum', 'Saum')"); err != nil {
			return err
		}
		for _, project := range projects {
			if _, err = tx.Exec("INSERT INTO projects_next (id, factory_id, name) VALUES (?, 'saum', ?)", project.ID, project.Name); err != nil {
				return err
			}
		}
	}
	if _, err = tx.Exec("DROP TABLE IF EXISTS projects; ALTER TABLE projects_next RENAME TO projects; PRAGMA user_version = 1;"); err != nil {
		return err
	}
	return tx.Commit()
}

func previousProjects(tx *sql.Tx, legacyFile string) ([]Project, error) {
	projects := []Project{}
	var exists int
	if err := tx.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'projects'").Scan(&exists); err != nil {
		return nil, err
	}
	if exists != 0 {
		rows, err := tx.Query("SELECT id, name FROM projects ORDER BY rowid")
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var project Project
			if err := rows.Scan(&project.ID, &project.Name); err != nil {
				_ = rows.Close()
				return nil, err
			}
			projects = append(projects, project)
		}
		err = errors.Join(rows.Err(), rows.Close())
		if err != nil {
			return nil, err
		}
	}
	bytes, err := os.ReadFile(legacyFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err == nil {
		var legacy []Project
		if err = json.Unmarshal(bytes, &legacy); err != nil {
			return nil, fmt.Errorf("read previous projects: %w", err)
		}
		projects = append(projects, legacy...)
	}
	unique := []Project{}
	seen := map[string]string{}
	for _, project := range projects {
		name, err := parseName(project.Name)
		if err != nil || name != project.Name || project.ID == "" {
			return nil, fmt.Errorf("invalid previous project")
		}
		if previous, ok := seen[project.ID]; ok {
			if previous != project.Name {
				return nil, fmt.Errorf("conflicting previous project %q", project.ID)
			}
			continue
		}
		seen[project.ID] = project.Name
		unique = append(unique, project)
	}
	return unique, nil
}
