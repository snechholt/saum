package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Settings struct {
	FactoriesRoot string `json:"factoriesRoot"`
	Theme         string `json:"theme"`
}

func migrateSettings(db *sql.DB, dataDir string) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version == 8 {
		return nil
	}
	root, err := filepath.Abs(filepath.Join(dataDir, "factories"))
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`CREATE TABLE settings(id INTEGER PRIMARY KEY CHECK(id=1),factories_root TEXT NOT NULL,theme TEXT NOT NULL CHECK(theme IN ('system','light','dark')));
 ALTER TABLE factories ADD COLUMN workspace_root TEXT NOT NULL DEFAULT '';
 ALTER TABLE factories ADD COLUMN repository_path TEXT NOT NULL DEFAULT '';`); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO settings VALUES (1,?,'system')", root); err != nil {
		return err
	}
	rows, err := tx.Query("SELECT id FROM factories")
	if err != nil {
		return err
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, id := range ids {
		workspace := filepath.Join(root, id)
		if _, err := tx.Exec("UPDATE factories SET workspace_root=?,repository_path=? WHERE id=?", workspace, filepath.Join(workspace, "repo"), id); err != nil {
			return err
		}
	}
	if _, err := tx.Exec("PRAGMA user_version=8"); err != nil {
		return err
	}
	return tx.Commit()
}
func parseSettings(input Settings) (Settings, error) {
	path := strings.TrimSpace(input.FactoriesRoot)
	if !utf8.ValidString(path) || strings.ContainsAny(path, "\x00\r\n") || len(path) > 4000 || !filepath.IsAbs(path) {
		return Settings{}, fmt.Errorf("factories directory must be an absolute path")
	}
	path = filepath.Clean(path)
	if filepath.Dir(path) == path {
		return Settings{}, fmt.Errorf("choose a factories directory, not the filesystem root")
	}
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return Settings{}, fmt.Errorf("factories directory points to a file")
	} else if err != nil && !os.IsNotExist(err) {
		return Settings{}, fmt.Errorf("cannot access factories directory: %w", err)
	}
	if input.Theme != "system" && input.Theme != "light" && input.Theme != "dark" {
		return Settings{}, fmt.Errorf("theme must be system, light, or dark")
	}
	return Settings{FactoriesRoot: path, Theme: input.Theme}, nil
}
func (s *Store) settings() (Settings, error) {
	var settings Settings
	err := s.db.QueryRow("SELECT factories_root,theme FROM settings WHERE id=1").Scan(&settings.FactoriesRoot, &settings.Theme)
	return settings, err
}
func (s *Store) saveSettings(input Settings) (Settings, error) {
	settings, err := parseSettings(input)
	if err != nil {
		return Settings{}, err
	}
	_, err = s.db.Exec("UPDATE settings SET factories_root=?,theme=? WHERE id=1", settings.FactoriesRoot, settings.Theme)
	return settings, err
}
func (s *Store) settingsHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		settings, err := s.settings()
		if err != nil {
			problem(w, 500, "Could not load settings")
			return
		}
		jsonResponse(w, 200, settings)
		return
	}
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/json" {
		problem(w, 415, "Expected application/json")
		return
	}
	var input Settings
	r.Body = http.MaxBytesReader(w, r.Body, 20000)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		problem(w, 400, "Invalid settings")
		return
	}
	settings, err := s.saveSettings(input)
	if err != nil {
		problem(w, 400, err.Error())
		return
	}
	jsonResponse(w, 200, settings)
}
func (s *Store) factoryWorkspace(id string) (workspace, repository string, err error) {
	err = s.db.QueryRow("SELECT workspace_root,repository_path FROM factories WHERE id=?", id).Scan(&workspace, &repository)
	return
}
