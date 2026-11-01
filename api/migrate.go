package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type previousProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Preserve older entries and their IDs during the terminology cutover. Keep source JSON untouched.
func migrate(db *sql.DB, legacyFile string) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version == 6 {
		return nil
	}
	if version == 5 {
		return migrateImplementations(db)
	}
	if version == 4 {
		return migratePlanReview(db)
	}
	if version == 3 {
		return migrateSessions(db)
	}
	if version < 0 || version > 6 {
		return fmt.Errorf("unsupported database version %d", version)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if version == 0 {
		projects, err := previousProjects(tx, legacyFile)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(`CREATE TABLE factories (id TEXT PRIMARY KEY NOT NULL, name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 80));`); err != nil {
			return err
		}
		if _, err = tx.Exec(storySchema); err != nil {
			return err
		}
		if len(projects) > 0 {
			if _, err = tx.Exec("INSERT INTO factories (id, name) VALUES ('saum', 'Saum')"); err != nil {
				return err
			}
			for _, project := range projects {
				if _, err = tx.Exec("INSERT INTO user_stories (id, factory_id, title, goal) VALUES (?, 'saum', ?, '')", project.ID, project.Name); err != nil {
					return err
				}
			}
		}
	} else {
		if _, err = tx.Exec(storySchema); err != nil {
			return err
		}
		statement := "INSERT INTO user_stories (id, factory_id, title, goal) SELECT id, factory_id, name, '' FROM projects ORDER BY rowid"
		if version == 2 {
			statement = "INSERT INTO user_stories (id, factory_id, title, goal) SELECT id, factory_id, name, goal FROM projects ORDER BY rowid"
		}
		if _, err = tx.Exec(statement); err != nil {
			return err
		}
	}
	if _, err = tx.Exec("DROP TABLE IF EXISTS projects; PRAGMA user_version = 3;"); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return migrateSessions(db)
}

const storySchema = `CREATE TABLE user_stories (
 id TEXT PRIMARY KEY NOT NULL,
 factory_id TEXT NOT NULL REFERENCES factories(id),
 title TEXT NOT NULL CHECK (length(trim(title)) BETWEEN 1 AND 80),
 goal TEXT NOT NULL DEFAULT '' CHECK (length(goal) <= 4000)
);`

func previousProjects(tx *sql.Tx, legacyFile string) ([]previousProject, error) {
	projects := []previousProject{}
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
			var project previousProject
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
		var legacy []previousProject
		if err = json.Unmarshal(bytes, &legacy); err != nil {
			return nil, fmt.Errorf("read previous projects: %w", err)
		}
		projects = append(projects, legacy...)
	}
	unique := []previousProject{}
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
