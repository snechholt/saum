package main

import (
	"errors"
	"net/http"
	"strconv"
)

type AuditSession struct {
	ID         string        `json:"id"`
	StoryID    string        `json:"storyId"`
	Status     SessionStatus `json:"status"`
	Kind       SessionKind   `json:"kind"`
	CreatedAt  string        `json:"createdAt"`
	FinishedAt string        `json:"finishedAt"`
	ThreadID   string        `json:"threadId"`
	Error      string        `json:"error"`
}
type AuditEntry struct {
	Cursor        int64        `json:"cursor"`
	Session       AuditSession `json:"session"`
	Factory       Factory      `json:"factory"`
	StoryTitle    string       `json:"storyTitle"`
	WorkItemTitle string       `json:"workItemTitle"`
}

func (s *Store) audit(before int64, factoryID string, kind SessionKind) ([]AuditEntry, error) {
	rows, err := s.db.Query(`SELECT run.rowid,run.id,run.story_id,run.status,run.kind,run.created_at,run.finished_at,run.thread_id,run.error,factory.id,factory.name,factory.workspace_root,factory.repository_path,story.title,COALESCE(item.title,'')
 FROM sessions AS run JOIN user_stories AS story ON story.id=run.story_id JOIN factories AS factory ON factory.id=story.factory_id
 LEFT JOIN work_item_sessions AS assigned ON assigned.session_id=run.id LEFT JOIN action_items AS item ON item.id=assigned.work_item_id
 WHERE (?=0 OR run.rowid<?) AND (?='' OR factory.id=?) AND (?='' OR run.kind=?) ORDER BY run.rowid DESC LIMIT 100`, before, before, factoryID, factoryID, kind, kind)
	if err != nil {
		return nil, err
	}
	entries := []AuditEntry{}
	for rows.Next() {
		var entry AuditEntry
		if err := rows.Scan(&entry.Cursor, &entry.Session.ID, &entry.Session.StoryID, &entry.Session.Status, &entry.Session.Kind, &entry.Session.CreatedAt, &entry.Session.FinishedAt, &entry.Session.ThreadID, &entry.Session.Error, &entry.Factory.ID, &entry.Factory.Name, &entry.Factory.WorkspaceRoot, &entry.Factory.RepositoryPath, &entry.StoryTitle, &entry.WorkItemTitle); err != nil {
			_ = rows.Close()
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return nil, err
	}
	return entries, nil
}
func (s *Store) auditHTTP(w http.ResponseWriter, r *http.Request) {
	before := int64(0)
	if raw := r.URL.Query().Get("before"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || value < 1 {
			problem(w, 400, "Invalid audit cursor")
			return
		}
		before = value
	}
	kind := SessionKind(r.URL.Query().Get("role"))
	if kind != "" && kind != PlanningSession && kind != ImplementationSession && kind != SummarySession {
		problem(w, 400, "Invalid agent role")
		return
	}
	entries, err := s.audit(before, r.URL.Query().Get("factory"), kind)
	if err != nil {
		problem(w, 500, "Could not load audit sessions")
		return
	}
	cursor := int64(0)
	if len(entries) == 100 {
		cursor = entries[len(entries)-1].Cursor
	}
	jsonResponse(w, 200, struct {
		Entries    []AuditEntry `json:"entries"`
		NextCursor int64        `json:"nextCursor"`
	}{entries, cursor})
}
