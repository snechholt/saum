package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestApprovalQueuesExactlyOneImplementation(t *testing.T) {
	store := testStore(t)
	plan, item := submittedPlan(t, store)
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	build, err := store.latestImplementation(plan.ID)
	if err != nil || build.Kind != ImplementationSession || build.ParentSessionID != plan.ID {
		t.Fatal(build, err)
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	var count int
	_ = store.db.QueryRow("SELECT count(*) FROM implementations WHERE plan_session_id=?", plan.ID).Scan(&count)
	if count != 1 {
		t.Fatal("approval scheduled duplicate work", count)
	}
	details, err := store.implementation(build.ID)
	if err != nil || len(details.Items) != 1 || details.Items[0].Key != item.Key {
		t.Fatal(details, err)
	}
	build, err = store.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := store.startSession(build.ID, plan.Worktree, "implement", "token"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.submitItem(build.ID, ItemSubmission{Key: "invented", Title: "Other", Description: "Other", AcceptanceCriteria: []string{"Other"}}); err != ErrPlanClosed {
		t.Fatal("implementation could rewrite plan", err)
	}
	if err := store.completePlan(build.ID); err != ErrPlanClosed {
		t.Fatal("implementation could submit another plan", err)
	}
	if err := store.completeImplementation(build.ID, "Done"); err == nil {
		t.Fatal("incomplete work accepted")
	}
	if err := store.reportImplementationItem(build.ID, ItemProgress{Key: item.Key, Summary: "Done"}, false); err == nil {
		t.Fatal("completed unstarted work")
	}
	if _, err := store.startImplementationItem(build.ID, item.Key); err != nil {
		t.Fatal(err)
	}
	if err := store.reportImplementationItem(build.ID, ItemProgress{Key: item.Key, Summary: "Need the recovery channel."}, true); err != nil {
		t.Fatal(err)
	}
	if err := store.completeImplementation(build.ID, "Done"); err == nil {
		t.Fatal("blocked work accepted")
	}
	if _, err := store.startImplementationItem(build.ID, item.Key); err != nil {
		t.Fatal(err)
	}
	if err := store.reportImplementationItem(build.ID, ItemProgress{Key: item.Key, Summary: "Implemented and tested."}, false); err != nil {
		t.Fatal(err)
	}
	if err := store.completeImplementation(build.ID, "All checks pass."); err != nil {
		t.Fatal(err)
	}
	if _, err := store.startImplementationItem(build.ID, item.Key); err == nil {
		t.Fatal("sealed implementation changed")
	}
	if err := store.finishSession(build.ID, SessionFailed, "test interruption"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.retryStory(plan.StoryID); err == nil {
		t.Fatal("approved story reverted to planning")
	}
	retry, err := store.retryImplementation(plan.StoryID)
	if err != nil || retry.Kind != ImplementationSession || retry.ParentSessionID != plan.ID {
		t.Fatal(retry, err)
	}
	original, _ := store.items(plan.ID)
	if len(original) != 1 || original[0].Title != item.Title {
		t.Fatal("approved snapshot changed")
	}
}
func TestExistingApprovalIsCaughtUpOnlyOnce(t *testing.T) {
	store := testStore(t)
	plan, _ := submittedPlan(t, store)
	if _, err := store.db.Exec("UPDATE sessions SET review_status='approved' WHERE id=?", plan.ID); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := store.queueApprovedPlans(); err != nil {
			t.Fatal(err)
		}
	}
	build, err := store.latestImplementation(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(build.ID, SessionCancelled, "stopped"); err != nil {
		t.Fatal(err)
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	latest, err := store.latestImplementation(plan.ID)
	if err != nil || latest.ID != build.ID {
		t.Fatal("cancelled build automatically restarted", latest, err)
	}
}
func TestImplementationHelperProcess(t *testing.T) {
	if os.Getenv("SAUM_TEST_IMPLEMENTATION") != "1" {
		return
	}
	_, _ = io.ReadAll(os.Stdin)
	fmt.Println(`{"type":"thread.started","thread_id":"implementation-thread"}`)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "implementation-test", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: os.Getenv("SAUM_TEST_MCP_URL"), HTTPClient: &http.Client{Transport: bearerTransport{os.Getenv("SAUM_SESSION_TOKEN")}}, DisableStandaloneSSE: true}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	call := func(code string) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "execute", Arguments: map[string]string{"code": code}})
		if err != nil || result.IsError {
			fmt.Fprintln(os.Stderr, result, err)
			os.Exit(3)
		}
	}
	call(`return saum.start_item({key:"login"});`)
	if err := os.WriteFile("main.go", []byte("package main\nimport \"fmt\"\nfunc main(){fmt.Print(\"implemented\")}\n"), 0600); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
	output, err := exec.Command("go", "run", "main.go").CombinedOutput()
	if err != nil || string(output) != "implemented" {
		fmt.Fprintln(os.Stderr, string(output), err)
		os.Exit(5)
	}
	call(`return saum.complete_item({key:"login",summary:"Created runnable software; go run main.go prints implemented."});`)
	if os.Getenv("SAUM_TEST_SKIP_COMPLETION") != "1" {
		call(`return saum.complete_implementation({summary:"Verified with go run main.go."});`)
	}
	_ = session.Close()
	fmt.Println(`{"type":"turn.completed"}`)
	os.Exit(0)
}
func TestApprovedPlanRunsSoftwareImplementation(t *testing.T) {
	for _, skip := range []bool{false, true} {
		t.Run(fmt.Sprint("skip-completion=", skip), func(t *testing.T) {
			store := testStore(t)
			plan, _ := submittedPlan(t, store)
			if err := store.approvePlan(plan.ID); err != nil {
				t.Fatal(err)
			}
			build, err := store.claimSession()
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(store.mcpHandler())
			defer server.Close()
			t.Setenv("SAUM_TEST_MCP_URL", server.URL)
			t.Setenv("SAUM_TEST_IMPLEMENTATION", "1")
			if skip {
				t.Setenv("SAUM_TEST_SKIP_COMPLETION", "1")
			}
			binary, _ := os.Executable()
			script := filepath.Join(t.TempDir(), "codex")
			if err := os.WriteFile(script, []byte(fmt.Sprintf("#!/bin/sh\nexec %q -test.run=^TestImplementationHelperProcess$ -- \"$@\"\n", binary)), 0700); err != nil {
				t.Fatal(err)
			}
			planner := Planner{Store: store, Root: filepath.Dir(plan.Worktree), Codex: script, MCPURL: server.URL, Timeout: 30 * time.Second}
			planner.execute(context.Background(), build)
			result, err := store.session(build.ID)
			if err != nil {
				t.Fatal(err)
			}
			expected := SessionCompleted
			if skip {
				expected = SessionFailed
			}
			if result.Status != expected {
				events, _ := store.events(build.ID, 0)
				data, _ := json.Marshal(events)
				t.Fatal(result, string(data))
			}
			if !strings.Contains(result.Prompt, "APPROVED") || strings.Contains(result.Prompt, "Only decompose") {
				t.Fatal("implementation received planning instructions")
			}
			output, err := exec.Command("go", "run", filepath.Join(result.Worktree, "main.go")).CombinedOutput()
			if err != nil || string(output) != "implemented" {
				t.Fatal("software artifact did not run", string(output), err)
			}
			detail, err := store.implementation(build.ID)
			if err != nil || detail.Items[0].Status != "completed" {
				t.Fatal(detail, err)
			}
		})
	}
}
