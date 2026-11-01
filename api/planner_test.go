package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type bearerTransport struct{ token string }

func (b bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	copy := r.Clone(r.Context())
	copy.Header.Set("Authorization", "Bearer "+b.token)
	return http.DefaultTransport.RoundTrip(copy)
}
func testQueuedSession(t *testing.T, store *Store) Session {
	t.Helper()
	factory, err := store.createFactory("Planner factory")
	if err != nil {
		t.Fatal(err)
	}
	spec, err := parseStorySpec("Sign in", "People can sign in and access only their private information.")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.createStory(factory.ID, spec); err != nil {
		t.Fatal(err)
	}
	if err = store.enqueueStories(); err != nil {
		t.Fatal(err)
	}
	session, err := store.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	return session
}
func fakeCodex(t *testing.T, mode string) string {
	t.Helper()
	t.Setenv("SAUM_TEST_CODEX_MODE", mode)
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "codex")
	script := fmt.Sprintf("#!/bin/sh\nexec %q -test.run=^TestCodexHelperProcess$ -- \"$@\"\n", binary)
	if err := os.WriteFile(file, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	return file
}
func TestCodexHelperProcess(t *testing.T) {
	mode := os.Getenv("SAUM_TEST_CODEX_MODE")
	if mode == "" {
		return
	}
	if _, err := io.ReadAll(os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(`{"type":"thread.started","thread_id":"test-thread"}`)
	fmt.Println(`{"type":"new.event.kind","text":"unrecognized events must survive 🎯"}`)
	if mode == "hang" {
		time.Sleep(30 * time.Second)
		os.Exit(0)
	}
	if mode == "failure" {
		fmt.Println(`{"type":"turn.failed","error":{"message":"test provider failure"}}`)
		os.Exit(1)
	}
	if mode == "full" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client := mcp.NewClient(&mcp.Implementation{Name: "codex-test", Version: "1"}, nil)
		session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: os.Getenv("SAUM_TEST_MCP_URL"), HTTPClient: &http.Client{Transport: bearerTransport{os.Getenv("SAUM_SESSION_TOKEN")}}, DisableStandaloneSSE: true}, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		input := ItemSubmission{Key: "authentication", Title: "Implement sign-in", Description: "Deliver sign-in for existing users without exposing private information.", AcceptanceCriteria: []string{"An authenticated user can access their own private information", "An unauthenticated request is rejected"}}
		for range 2 {
			result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "saum_item_submit", Arguments: input})
			if err != nil || result.IsError {
				fmt.Fprintf(os.Stderr, "submit failed: %+v %v\n", result, err)
				os.Exit(3)
			}
		}
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "saum_plan_complete", Arguments: struct{}{}})
		if err != nil || result.IsError {
			fmt.Fprintf(os.Stderr, "complete failed: %+v %v\n", result, err)
			os.Exit(4)
		}
		_ = session.Close()
		if err = os.WriteFile("plan.txt", []byte("Planning artifact retained in the session worktree."), 0600); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(5)
		}
	}
	fmt.Println(`{"type":"item.completed","item":{"type":"agent_message","text":"Plan submitted."}}`)
	fmt.Println(`{"type":"turn.completed","usage":{"input_tokens":12,"output_tokens":34}}`)
	os.Exit(0)
}
func TestPlannerPersistsRealMCPCallsAndWorktreeArtifacts(t *testing.T) {
	store := testStore(t)
	session := testQueuedSession(t, store)
	server := httptest.NewServer(store.mcpHandler())
	defer server.Close()
	t.Setenv("SAUM_TEST_MCP_URL", server.URL)
	planner := Planner{Store: store, Root: t.TempDir(), Codex: fakeCodex(t, "full"), MCPURL: server.URL, Timeout: 20 * time.Second}
	planner.execute(context.Background(), session)
	result, err := store.session(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != SessionCompleted || !result.PlanSubmitted || result.ThreadID != "test-thread" {
		events, _ := store.events(session.ID, 0)
		t.Fatalf("session failed: %+v %v; events=%+v", result, err, events)
	}
	items, err := store.items(session.ID)
	if err != nil || len(items) != 1 || items[0].Key != "authentication" {
		t.Fatalf("missing or duplicated items: %+v %v", items, err)
	}
	if _, err := os.Stat(filepath.Join(result.Worktree, "plan.txt")); err != nil {
		t.Fatal("worktree artifact missing", err)
	}
	events, err := store.events(session.ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	stdout := ""
	tools := 0
	for _, event := range events {
		if event.Stream == "stdout" {
			stdout += string(event.Data)
		}
		if event.Stream == "tool" {
			tools++
		}
	}
	if !strings.Contains(stdout, "unrecognized events must survive 🎯") || !strings.Contains(stdout, "input_tokens") || tools != 3 {
		t.Fatalf("session output not retained: %q tools=%d", stdout, tools)
	}
	if err := store.enqueueStories(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.claimSession(); err != sql.ErrNoRows {
		t.Fatal("completed story was automatically restarted", err)
	}
}
func TestProcessExitIsNotPlanCompletion(t *testing.T) {
	for _, mode := range []string{"no-tools", "failure"} {
		t.Run(mode, func(t *testing.T) {
			store := testStore(t)
			session := testQueuedSession(t, store)
			planner := Planner{Store: store, Root: t.TempDir(), Codex: fakeCodex(t, mode), MCPURL: "http://127.0.0.1:1/mcp"}
			planner.execute(context.Background(), session)
			result, err := store.session(session.ID)
			if err != nil || result.Status != SessionFailed || result.Error == "" {
				t.Fatalf("false completion: %+v %v", result, err)
			}
		})
	}
}
func TestPlanningCancellationAndRecovery(t *testing.T) {
	store := testStore(t)
	session := testQueuedSession(t, store)
	planner := Planner{Store: store, Root: t.TempDir(), Codex: fakeCodex(t, "hang"), MCPURL: "http://127.0.0.1:1/mcp"}
	done := make(chan struct{})
	go func() { planner.execute(context.Background(), session); close(done) }()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		current, _ := store.session(session.ID)
		if current.ThreadID != "" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !planner.Cancel(session.ID) {
		t.Fatal("running session not cancellable")
	}
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("cancel did not reap subprocess")
	}
	current, err := store.session(session.ID)
	if err != nil || current.Status != SessionCancelled {
		t.Fatalf("not cancelled: %+v %v", current, err)
	}
	retry, err := store.retryStory(session.StoryID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.retryStory(session.StoryID); err == nil {
		t.Fatal("accepted concurrent planning attempts")
	}
	claimed, err := store.claimSession()
	if err != nil || claimed.ID != retry.ID {
		t.Fatal("retry not queued", err)
	}
	if err = planner.recoverSessions(); err != nil {
		t.Fatal(err)
	}
	recovered, err := store.session(retry.ID)
	if err != nil || recovered.Status != SessionInterrupted {
		t.Fatalf("not recovered: %+v %v", recovered, err)
	}
}
func TestSessionToolsRequireScopedActiveToken(t *testing.T) {
	store := testStore(t)
	session := testQueuedSession(t, store)
	if err := store.startSession(session.ID, t.TempDir(), "prompt", "secret"); err != nil {
		t.Fatal(err)
	}
	handler := store.mcpHandler()
	request := httptest.NewRequest("POST", "http://127.0.0.1/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 401 {
		t.Fatal("MCP accepted an unauthenticated caller", response.Code)
	}
	if err := store.completePlan(session.ID); err == nil {
		t.Fatal("accepted an empty plan")
	}
	input := ItemSubmission{Key: "same", Title: "First", Description: "Work", AcceptanceCriteria: []string{"Works"}}
	if _, err := store.submitItem(session.ID, input); err != nil {
		t.Fatal(err)
	}
	input.Title = "Conflicting"
	if _, err := store.submitItem(session.ID, input); err == nil {
		t.Fatal("overwrote existing item")
	}
	if err := store.completePlan(session.ID); err != nil {
		t.Fatal(err)
	}
	input.Key = "new"
	if _, err := store.submitItem(session.ID, input); err != ErrPlanClosed {
		t.Fatal("accepted item after plan completion", err)
	}
	if err := store.finishSession(session.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := store.authenticatedSession("secret"); err == nil {
		t.Fatal("terminal session token not revoked")
	}
	events, err := store.events(session.ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		if event.Stream == "tool" && !json.Valid(event.Data) {
			t.Fatal("tool record is not structured JSON")
		}
	}
}
