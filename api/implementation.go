package main

import "database/sql"

type SessionKind string

const (
	PlanningSession       SessionKind = "planning"
	ImplementationSession SessionKind = "implementation"
	SummarySession        SessionKind = "summary"
)

func migrateImplementations(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`ALTER TABLE sessions ADD COLUMN kind TEXT NOT NULL DEFAULT 'planning' CHECK(kind IN ('planning','implementation'));
 CREATE TABLE implementations (
 session_id TEXT PRIMARY KEY NOT NULL REFERENCES sessions(id),plan_session_id TEXT NOT NULL REFERENCES sessions(id),
 submitted INTEGER NOT NULL DEFAULT 0,summary TEXT NOT NULL DEFAULT ''
 );
 CREATE INDEX builds_by_plan ON implementations(plan_session_id);
 CREATE TABLE implementation_items (
 session_id TEXT NOT NULL REFERENCES implementations(session_id),plan_item_id TEXT NOT NULL REFERENCES action_items(id),item_key TEXT NOT NULL,
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','running','completed','blocked')),
 summary TEXT NOT NULL DEFAULT '',PRIMARY KEY(session_id,item_key)
 );
 CREATE UNIQUE INDEX one_running_item ON implementation_items(session_id) WHERE status='running';
 PRAGMA user_version=6;`)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return migrateWorkItems(db)
}
