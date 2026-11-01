package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ItemSubmission struct {
	Key                string   `json:"key" jsonschema:"Stable unique key for this item; reuse exactly when retrying the same submission"`
	Title              string   `json:"title" jsonschema:"Short actionable item title"`
	Description        string   `json:"description" jsonschema:"Specific bounded work to perform and its relationship to the user story"`
	AcceptanceCriteria []string `json:"acceptanceCriteria" jsonschema:"Observable conditions that establish this item is complete"`
}

func parseItem(input ItemSubmission) (ItemSubmission, error) {
	input.Key = strings.TrimSpace(input.Key)
	if len(input.Key) < 1 || len(input.Key) > 80 {
		return input, fmt.Errorf("key must be 1–80 bytes")
	}
	title, err := parseName(input.Title)
	if err != nil {
		return input, err
	}
	input.Title = title
	input.Description = strings.TrimSpace(input.Description)
	if !utf8.ValidString(input.Description) || len(input.Description) < 1 || utf8.RuneCountInString(input.Description) > 8000 {
		return input, fmt.Errorf("description must be 1–8000 characters")
	}
	if len(input.AcceptanceCriteria) < 1 || len(input.AcceptanceCriteria) > 20 {
		return input, fmt.Errorf("provide 1–20 acceptance criteria")
	}
	for i, criterion := range input.AcceptanceCriteria {
		criterion = strings.TrimSpace(criterion)
		if !utf8.ValidString(criterion) || criterion == "" || utf8.RuneCountInString(criterion) > 2000 {
			return input, fmt.Errorf("each acceptance criterion must be 1–2000 characters")
		}
		input.AcceptanceCriteria[i] = criterion
	}
	return input, nil
}
func (s *Store) submitItem(sessionID string, raw ItemSubmission) (ActionItem, error) {
	input, err := parseItem(raw)
	if err != nil {
		return ActionItem{}, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return ActionItem{}, err
	}
	defer tx.Rollback()
	var status SessionStatus
	var closed bool
	if err = tx.QueryRow("SELECT status,plan_submitted FROM sessions WHERE id=?", sessionID).Scan(&status, &closed); err != nil {
		return ActionItem{}, err
	}
	if status != SessionRunning || closed {
		return ActionItem{}, ErrPlanClosed
	}
	criteria, err := json.Marshal(input.AcceptanceCriteria)
	if err != nil {
		return ActionItem{}, err
	}
	item := ActionItem{ID: rand.Text(), SessionID: sessionID, Key: input.Key, Title: input.Title, Description: input.Description, AcceptanceCriteria: input.AcceptanceCriteria}
	var previous ActionItem
	var previousCriteria string
	err = tx.QueryRow("SELECT id,title,description,acceptance_criteria FROM action_items WHERE session_id=? AND item_key=?", sessionID, input.Key).Scan(&previous.ID, &previous.Title, &previous.Description, &previousCriteria)
	if err == nil {
		if previous.Title != item.Title || previous.Description != item.Description || previousCriteria != string(criteria) {
			return ActionItem{}, fmt.Errorf("key already used for a different item")
		}
		item.ID = previous.ID
	} else if errors.Is(err, sql.ErrNoRows) {
		var count int
		if err = tx.QueryRow("SELECT count(*) FROM action_items WHERE session_id=?", sessionID).Scan(&count); err != nil {
			return ActionItem{}, err
		}
		if count >= 100 {
			return ActionItem{}, fmt.Errorf("plan item limit reached (100)")
		}
		if _, err = tx.Exec("INSERT INTO action_items (id,session_id,item_key,title,description,acceptance_criteria) VALUES (?,?,?,?,?,?)", item.ID, sessionID, item.Key, item.Title, item.Description, string(criteria)); err != nil {
			return ActionItem{}, err
		}
	} else {
		return ActionItem{}, err
	}
	event, _ := json.Marshal(map[string]any{"tool": "saum.submit_item", "input": input, "result": item})
	if _, err = tx.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,'tool',?)", sessionID, event); err != nil {
		return ActionItem{}, err
	}
	if err = tx.Commit(); err != nil {
		return ActionItem{}, err
	}
	return item, nil
}
func (s *Store) completePlan(sessionID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var status SessionStatus
	if err = tx.QueryRow("SELECT status FROM sessions WHERE id=?", sessionID).Scan(&status); err != nil {
		return err
	}
	if status != SessionRunning {
		return ErrPlanClosed
	}
	var count int
	if err = tx.QueryRow("SELECT count(*) FROM action_items WHERE session_id=?", sessionID).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("submit at least one actionable item before completing the plan")
	}
	var unresolved int
	if err := tx.QueryRow(`SELECT count(*) FROM plan_comments WHERE batch_id=(SELECT feedback_batch_id FROM sessions WHERE id=?) AND NOT EXISTS (SELECT 1 FROM feedback_resolutions WHERE session_id=? AND comment_id=plan_comments.id)`, sessionID, sessionID).Scan(&unresolved); err != nil {
		return err
	}
	if unresolved > 0 {
		return fmt.Errorf("address all %d remaining feedback comments with saum.resolve_feedback before completing the plan", unresolved)
	}
	if _, err = tx.Exec("UPDATE sessions SET plan_submitted=1 WHERE id=?", sessionID); err != nil {
		return err
	}
	if _, err = tx.Exec("INSERT INTO session_events (session_id,stream,data) VALUES (?,'tool',?)", sessionID, []byte(`{"tool":"saum.complete_plan","result":"plan submitted"}`)); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) mcpHandler() http.Handler {
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		session, err := s.authenticatedSession(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if err != nil {
			return nil
		}
		closedWorld := false
		server := mcp.NewServer(&mcp.Implementation{Name: "saum", Version: "0.1.0"}, &mcp.ServerOptions{Instructions: "Use search to discover session-scoped Saum functions. Use execute to run synchronous JavaScript with the saum namespace; return a JSON value. There is no filesystem, network, import, or process API inside execute. Functions commit independently. Decompose the assigned story; do not implement it."})
		mcp.AddTool(server, &mcp.Tool{Name: "search", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &closedWorld}, Description: "Discover Saum functions, signatures, and examples by keywords. An empty query lists all functions."}, func(ctx context.Context, request *mcp.CallToolRequest, input struct {
			Query string `json:"query"`
		}) (*mcp.CallToolResult, map[string]any, error) {
			functions := searchFunctions(input.Query)
			data, _ := json.Marshal(map[string]any{"tool": "search", "query": input.Query, "functions": functions})
			if err := s.appendEvent(session.ID, "tool", data); err != nil {
				return nil, nil, err
			}
			return nil, map[string]any{"functions": functions}, nil
		})
		mcp.AddTool(server, &mcp.Tool{Name: "execute", Annotations: &mcp.ToolAnnotations{DestructiveHint: &closedWorld, OpenWorldHint: &closedWorld}, Description: "Run synchronous JavaScript using the functions discovered with search. Example: return saum.submit_item({...});. No await/import/require. Return a JSON value. Each call commits independently; reuse keys for retries."}, func(ctx context.Context, request *mcp.CallToolRequest, input struct {
			Code string `json:"code"`
		}) (*mcp.CallToolResult, map[string]any, error) {
			result, err := s.executeCode(ctx, session.ID, input.Code)
			return nil, map[string]any{"result": result}, err
		})
		return server
	}, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true, MaxRequestBodyBytes: 128 << 10})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			problem(w, 401, "Session token required")
			return
		}
		if _, err := s.authenticatedSession(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")); err != nil {
			problem(w, 401, "Invalid or expired session token")
			return
		}
		handler.ServeHTTP(w, r)
	})
}
