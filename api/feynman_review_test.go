//go:build feynman_review

package main

// Opt-in contract probes for architecture/feynman-spec.md. The contract tests
// intentionally fail on the reviewed implementation; observations document gaps
// without pretending those gaps are an already-agreed product policy.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func feynmanGoal(t *testing.T, s *Store, plan Session) Session {
	t.Helper()
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	run, err := s.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := s.startSession(run.ID, plan.Worktree, "assigned item", "probe-token"); err != nil {
		t.Fatal(err)
	}
	return run
}

func TestFeynmanApprovalMustReferToReviewedBytes(t *testing.T) {
	s := testStore(t)
	plan := twoItemPlan(t, s)
	goal := feynmanGoal(t, s, plan)
	artifact := filepath.Join(plan.Worktree, "answer.go")
	if err := os.WriteFile(artifact, []byte("package app\nfunc Answer() int { return 42 }\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, item := prepareSummary(t, s, goal)
	if err := os.WriteFile(artifact, []byte("package app\nfunc Answer() int { panic(\"changed after review\") }\n"), 0600); err != nil {
		t.Fatal(err)
	}
	next, err := s.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID}, true)
	if err == nil {
		t.Fatalf("approved changed source bytes and unlocked next item: %v", next != nil)
	}
}

func TestFeynmanCriterionEvidenceMustBeCorrectable(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	goal := feynmanGoal(t, s, plan)
	index := 0
	if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Initial check", Evidence: []string{"initial command: pass"}}); err != nil {
		t.Fatal(err)
	}
	err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Reverified after a correction", Evidence: []string{"expanded regression command: pass"}})
	if err != nil {
		t.Fatalf("cannot supersede earlier evidence after correcting the implementation: %v", err)
	}
}

func TestFeynmanWorkspaceMustResolveInsideFactory(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	if err := s.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	run, err := s.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Rename(plan.Worktree, plan.Worktree+".original"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, plan.Worktree); err != nil {
		t.Fatal(err)
	}
	factoryID, _, _, err := s.storyContext(plan.StoryID)
	if err != nil {
		t.Fatal(err)
	}
	planner := Planner{Store: s}
	workspace, err := planner.sessionWorkspace(t.Context(), run, factoryID)
	if err == nil {
		resolved, resolveErr := filepath.EvalSymlinks(workspace)
		if resolveErr != nil {
			t.Fatal(resolveErr)
		}
		t.Fatalf("accepted managed-looking path %q that resolves outside the factory to %q", workspace, resolved)
	}
}

func TestFeynmanNativeThreadIdentityMustFailClosed(t *testing.T) {
	for _, mode := range []string{"missing", "changed"} {
		t.Run(mode, func(t *testing.T) {
			s := testStore(t)
			plan, _ := submittedPlan(t, s)
			if err := s.approvePlan(plan.ID); err != nil {
				t.Fatal(err)
			}
			expected := "00000000-0000-4000-8000-000000000001"
			if _, err := s.db.Exec("UPDATE goal_workflows SET thread_id=? WHERE plan_session_id=?", expected, plan.ID); err != nil {
				t.Fatal(err)
			}
			run, err := s.claimSession()
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(s.mcpHandler())
			defer server.Close()
			t.Setenv("SAUM_FEYNMAN_THREAD", mode)
			t.Setenv("SAUM_TEST_MCP_URL", server.URL)
			binary, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			script := filepath.Join(t.TempDir(), "codex")
			if err := os.WriteFile(script, []byte(fmt.Sprintf("#!/bin/sh\nexec %q -test.run=^TestFeynmanCodexHelperProcess$ -- \"$@\"\n", binary)), 0700); err != nil {
				t.Fatal(err)
			}
			planner := Planner{Store: s, Codex: script, MCPURL: server.URL, Timeout: 20 * time.Second}
			planner.execute(context.Background(), run)
			finished, err := s.session(run.ID)
			if err != nil {
				t.Fatal(err)
			}
			thread, err := s.sessionThread(finished)
			if err != nil {
				t.Fatal(err)
			}
			if finished.Status == SessionCompleted || thread != expected {
				t.Fatalf("accepted unverified thread continuity: status=%s observed=%q registry=%q expected=%q", finished.Status, finished.ThreadID, thread, expected)
			}
		})
	}
}

func TestFeynmanCodexHelperProcess(t *testing.T) {
	mode := os.Getenv("SAUM_FEYNMAN_THREAD")
	if mode == "" {
		return
	}
	prompt, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(2)
	}
	if mode == "changed" {
		fmt.Println(`{"type":"thread.started","thread_id":"00000000-0000-4000-8000-000000000002"}`)
	}
	start, end := strings.Index(string(prompt), "<assigned_work_item>"), strings.Index(string(prompt), "</assigned_work_item>")
	var item WorkItem
	if start < 0 || end < 0 || json.Unmarshal(prompt[start+len("<assigned_work_item>"):end], &item) != nil {
		os.Exit(3)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "feynman-probe", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: os.Getenv("SAUM_TEST_MCP_URL"), HTTPClient: &http.Client{Transport: bearerTransport{os.Getenv("SAUM_SESSION_TOKEN")}}, DisableStandaloneSSE: true}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
	call := func(code string) {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "execute", Arguments: map[string]string{"code": code}})
		if err != nil || result.IsError {
			fmt.Fprintln(os.Stderr, result, err)
			os.Exit(5)
		}
	}
	for index := range item.AcceptanceCriteria {
		call(fmt.Sprintf(`return saum.work_item.submit_acceptance_criteria({criterionIndex:%d,summary:"Claimed verification",evidence:["claimed test passed"]});`, index))
	}
	call("return saum.work_item.request_review();")
	_ = session.Close()
	fmt.Println(`{"type":"turn.completed"}`)
	os.Exit(0)
}

func TestFeynmanObservationNarrativeIsNotProof(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	goal := feynmanGoal(t, s, plan)
	for index := range 2 {
		if err := s.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Trust me", Evidence: []string{"I say it passed"}}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.requestWorkItemReview(goal.ID); err != nil {
		t.Fatal(err)
	}
	t.Log("ASSURANCE GAP: an empty implementation with no executed checks can reach review using only free-form claims; human approval is still required")
}

func TestFeynmanObservationAcceptedWorkIsNotIntegrated(t *testing.T) {
	s := testStore(t)
	plan, _ := submittedPlan(t, s)
	goal := feynmanGoal(t, s, plan)
	if err := os.WriteFile(filepath.Join(plan.Worktree, "delivered.txt"), []byte("approved work"), 0600); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		command := exec.Command("git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgsign=false", "-c", "user.name=Review probe", "-c", "user.email=probe@localhost", "-C", plan.Worktree}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatal(string(output), err)
		}
	}
	git("add", "delivered.txt")
	git("commit", "-m", "test: committed work in goal branch")
	_, item := prepareSummary(t, s, goal)
	if _, err := s.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID}, true); err != nil {
		t.Fatal(err)
	}
	factoryID, _, _, _ := s.storyContext(plan.StoryID)
	stories, err := s.listStories(factoryID)
	if err != nil || stories[0].GoalStatus != "completed" {
		t.Fatal(stories, err)
	}
	planner := Planner{Store: s}
	next, err := planner.worktree(t.Context(), factoryID, "another-goal")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(next, "delivered.txt")); !os.IsNotExist(err) {
		t.Fatalf("integration behavior changed; revisit this observation: %v", err)
	}
	t.Log("DELIVERY GAP: goal is completed, but a new goal created from main cannot see the approved committed file")
}
