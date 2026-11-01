package main

import "database/sql"

func migrateReviewSafety(db *sql.DB) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version == 9 {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`ALTER TABLE acceptance_submissions ADD COLUMN artifact_tree TEXT NOT NULL DEFAULT '';
 ALTER TABLE acceptance_submissions ADD COLUMN check_ids TEXT NOT NULL DEFAULT '[]';
 CREATE TABLE acceptance_history (
 id INTEGER PRIMARY KEY AUTOINCREMENT,session_id TEXT NOT NULL REFERENCES sessions(id),work_item_id TEXT NOT NULL REFERENCES work_items(id),criterion_index INTEGER NOT NULL,
 summary TEXT NOT NULL,evidence TEXT NOT NULL,artifact_tree TEXT NOT NULL,check_ids TEXT NOT NULL,created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ','now'))
 );
 INSERT INTO acceptance_history(session_id,work_item_id,criterion_index,summary,evidence,artifact_tree,check_ids) SELECT session_id,work_item_id,criterion_index,summary,evidence,'','[]' FROM acceptance_submissions;
 ALTER TABLE work_items ADD COLUMN review_tree TEXT NOT NULL DEFAULT '';
 ALTER TABLE work_item_reviews ADD COLUMN artifact_tree TEXT NOT NULL DEFAULT '';
 ALTER TABLE work_item_reviews ADD COLUMN approved_commit TEXT NOT NULL DEFAULT '';
 ALTER TABLE goal_workflows ADD COLUMN base_commit TEXT NOT NULL DEFAULT '';
 ALTER TABLE goal_workflows ADD COLUMN integrated_commit TEXT NOT NULL DEFAULT '';
 ALTER TABLE goal_workflows ADD COLUMN base_ref TEXT NOT NULL DEFAULT '';
 CREATE TABLE integrations(review_id TEXT PRIMARY KEY NOT NULL REFERENCES work_item_reviews(id),plan_session_id TEXT NOT NULL REFERENCES goal_workflows(plan_session_id),base_commit TEXT NOT NULL,base_ref TEXT NOT NULL,commit_id TEXT NOT NULL,status TEXT NOT NULL CHECK(status IN ('queued','running','completed','failed')),error TEXT NOT NULL DEFAULT '');
 CREATE TABLE verification_checks (
 id TEXT PRIMARY KEY NOT NULL,session_id TEXT NOT NULL REFERENCES sessions(id),work_item_id TEXT NOT NULL REFERENCES work_items(id),artifact_tree TEXT NOT NULL,
 command TEXT NOT NULL,status TEXT NOT NULL CHECK(status IN ('queued','running','passed','failed','cancelled')),
 output TEXT NOT NULL DEFAULT '',exit_code INTEGER,error TEXT NOT NULL DEFAULT '',process_id INTEGER NOT NULL DEFAULT 0,process_identity TEXT NOT NULL DEFAULT '',created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ','now')),finished_at TEXT NOT NULL DEFAULT ''
 );
 CREATE INDEX checks_by_session ON verification_checks(session_id);
 PRAGMA user_version=9;`)
	if err != nil {
		return err
	}
	return tx.Commit()
}
