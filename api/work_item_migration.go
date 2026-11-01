package main

import (
	"database/sql"
	"fmt"
)

func migrateWorkItems(db *sql.DB) error {
	// Rebuild the session kind constraint without changing IDs, row order, or any history.
	if _, err := db.Exec("PRAGMA foreign_keys=OFF"); err != nil {
		return err
	}
	defer db.Exec("PRAGMA foreign_keys=ON")
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`CREATE TABLE sessions_next (
 id TEXT PRIMARY KEY NOT NULL,story_id TEXT NOT NULL REFERENCES user_stories(id),
 status TEXT NOT NULL CHECK(status IN ('queued','preparing','running','completed','failed','interrupted','cancelled')),
 worktree TEXT NOT NULL DEFAULT '',thread_id TEXT NOT NULL DEFAULT '',prompt TEXT NOT NULL DEFAULT '',error TEXT NOT NULL DEFAULT '',
 plan_submitted INTEGER NOT NULL DEFAULT 0,token_hash BLOB,process_id INTEGER NOT NULL DEFAULT 0,process_identity TEXT NOT NULL DEFAULT '',
 created_at TEXT NOT NULL,finished_at TEXT NOT NULL DEFAULT '',parent_session_id TEXT REFERENCES sessions(id),feedback_batch_id TEXT REFERENCES plan_feedback_batches(id),
 review_status TEXT NOT NULL DEFAULT 'pending' CHECK(review_status IN ('pending','approved','changes_requested')),reviewed_at TEXT NOT NULL DEFAULT '',
 kind TEXT NOT NULL DEFAULT 'planning' CHECK(kind IN ('planning','implementation','summary'))
 );
 INSERT INTO sessions_next SELECT * FROM sessions ORDER BY rowid;
 DROP TABLE sessions;
 ALTER TABLE sessions_next RENAME TO sessions;
 CREATE UNIQUE INDEX one_active_session_per_story ON sessions(story_id) WHERE status IN ('queued','preparing','running');
 CREATE INDEX sessions_by_story ON sessions(story_id);
 CREATE TABLE goal_workflows(plan_session_id TEXT PRIMARY KEY NOT NULL REFERENCES sessions(id),thread_id TEXT NOT NULL DEFAULT '');
 CREATE TABLE factory_agents(factory_id TEXT PRIMARY KEY NOT NULL REFERENCES factories(id),thread_id TEXT NOT NULL DEFAULT '');
 CREATE TABLE work_items (
 id TEXT PRIMARY KEY NOT NULL REFERENCES action_items(id),plan_session_id TEXT NOT NULL REFERENCES goal_workflows(plan_session_id),position INTEGER NOT NULL,
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','queued','working','blocked','review_requested','summarizing','awaiting_review','approved','failed')),
 blocker TEXT NOT NULL DEFAULT '',current_session_id TEXT REFERENCES sessions(id),last_goal_session_id TEXT REFERENCES sessions(id),
 UNIQUE(plan_session_id,position)
 );
 CREATE TABLE work_item_sessions(session_id TEXT PRIMARY KEY NOT NULL REFERENCES sessions(id),work_item_id TEXT NOT NULL REFERENCES work_items(id));
 CREATE TABLE work_item_messages(id INTEGER PRIMARY KEY AUTOINCREMENT,work_item_id TEXT NOT NULL REFERENCES work_items(id),body TEXT NOT NULL);
 CREATE TABLE acceptance_submissions (
 session_id TEXT NOT NULL REFERENCES sessions(id),work_item_id TEXT NOT NULL REFERENCES work_items(id),criterion_index INTEGER NOT NULL CHECK(criterion_index>=0),
 summary TEXT NOT NULL,evidence TEXT NOT NULL,PRIMARY KEY(session_id,criterion_index)
 );
 CREATE TABLE work_item_reviews (
 id TEXT PRIMARY KEY NOT NULL,work_item_id TEXT NOT NULL REFERENCES work_items(id),goal_session_id TEXT NOT NULL REFERENCES sessions(id),summary_session_id TEXT NOT NULL UNIQUE REFERENCES sessions(id),
 summary TEXT NOT NULL,verification TEXT NOT NULL,risks TEXT NOT NULL,decision TEXT NOT NULL DEFAULT 'pending' CHECK(decision IN ('pending','approved','changes_requested')),
 feedback TEXT NOT NULL DEFAULT '',created_at TEXT NOT NULL DEFAULT(strftime('%Y-%m-%dT%H:%M:%fZ','now')),reviewed_at TEXT NOT NULL DEFAULT ''
 );
 INSERT INTO goal_workflows(plan_session_id) SELECT DISTINCT plan_session_id FROM implementations;
 INSERT INTO work_items(id,plan_session_id,position,status,blocker)
 SELECT item.id,item.session_id,(SELECT count(*)-1 FROM action_items WHERE session_id=item.session_id AND rowid<=item.rowid),
 CASE WHEN progress.status='blocked' THEN 'blocked' ELSE 'pending' END,CASE WHEN progress.status='blocked' THEN progress.summary ELSE '' END
 FROM action_items AS item JOIN goal_workflows AS flow ON flow.plan_session_id=item.session_id
 LEFT JOIN implementation_items AS progress ON progress.plan_item_id=item.id AND progress.session_id=(SELECT session_id FROM implementations WHERE plan_session_id=item.session_id ORDER BY rowid DESC LIMIT 1);
 PRAGMA user_version=7;`)
	if err != nil {
		return err
	}
	check, err := tx.Query("PRAGMA foreign_key_check")
	if err != nil {
		return err
	}
	broken := check.Next()
	if err := check.Close(); err != nil {
		return err
	}
	if broken {
		return fmt.Errorf("session migration broke a foreign-key reference")
	}
	return tx.Commit()
}
