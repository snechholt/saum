package main

import (
	"context"
	"database/sql"
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

func twoItemPlan(t *testing.T, store *Store) Session {
	t.Helper()
	run := testQueuedSession(t, store)
	if err := store.startSession(run.ID, t.TempDir(), "plan", "token"); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"first", "second"} {
		if _, err := store.submitItem(run.ID, ItemSubmission{Key: key, Title: key + " unique item", Description: "Deliver " + key, AcceptanceCriteria: []string{"Behavior works", "Tests pass"}}); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.completePlan(run.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(run.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	run, _ = store.session(run.ID)
	return run
}
func startClaimed(t *testing.T, s *Store) Session {
	t.Helper()
	run, err := s.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := s.startSession(run.ID, t.TempDir(), "prompt", "token"); err != nil {
		t.Fatal(err)
	}
	return run
}
func submitCriteria(t *testing.T, s *Store, run Session) {
	t.Helper()
	item, err := s.assignedWorkItem(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	for index := range item.AcceptanceCriteria {
		if err := s.submitAcceptanceCriterion(run.ID, CriterionInput{CriterionIndex: &index, Summary: "Verified behavior", Evidence: []string{"Automated test passed"}}); err != nil {
			t.Fatal(err)
		}
	}
}
func prepareSummary(t *testing.T, s *Store, goal Session) (Session, WorkItem) {
	t.Helper()
	submitCriteria(t, s, goal)
	if err := s.requestWorkItemReview(goal.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.finishSession(goal.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	if err := s.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	summary := startClaimed(t, s)
	if summary.Kind != SummarySession {
		t.Fatal(summary)
	}
	if err := s.submitExecutiveSummary(summary.ID, ExecutiveSummaryInput{Summary: "Delivered behavior.", Verification: "Observed passing tests.", Risks: []string{}}); err != nil {
		t.Fatal(err)
	}
	if err := s.finishSession(summary.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	item, err := s.assignedWorkItem(summary.ID)
	if err != nil {
		t.Fatal(err)
	}
	return summary, item
}
func TestHumanExecutiveApprovalIsTheOnlyNextItemGate(t *testing.T) {
	store := testStore(t)
	plan := twoItemPlan(t, store)
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, store)
	if goal.Kind != ImplementationSession {
		t.Fatal(goal)
	}
	if err := store.saveSessionThread(goal, "goal-thread"); err != nil {
		t.Fatal(err)
	}
	prompt, err := store.workItemPrompt(goal, "Story", "Goal")
	if err != nil || strings.Contains(prompt, "second unique item") {
		t.Fatal("goal received more than one item", err)
	}
	if err := store.requestWorkItemReview(goal.ID); err == nil {
		t.Fatal("requested review before criteria")
	}
	index := 1
	if err := store.submitAcceptanceCriterion(goal.ID, CriterionInput{CriterionIndex: &index, Summary: "Verified", Evidence: []string{"pass"}}); err == nil {
		t.Fatal("criteria submitted out of order")
	}
	if _, err := store.submitItem(goal.ID, ItemSubmission{Key: "new", Title: "New", Description: "New", AcceptanceCriteria: []string{"New"}}); err != ErrPlanClosed {
		t.Fatal("goal rewrote approved plan", err)
	}
	summary, item := prepareSummary(t, store, goal)
	if err := store.saveSessionThread(summary, "factory-thread"); err != nil {
		t.Fatal(err)
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.claimSession(); err != sql.ErrNoRows {
		t.Fatal("next item started before human approval", err)
	}
	review := item.Reviews[0]
	feedback, err := store.decideWorkItem(item.ID, WorkItemDecision{ReviewID: review.ID, Feedback: "Tighten the error handling."}, false)
	if err != nil || feedback == nil {
		t.Fatal(feedback, err)
	}
	revised := startClaimed(t, store)
	assigned, _ := store.assignedWorkItem(revised.ID)
	if assigned.ID != item.ID {
		t.Fatal("feedback advanced to next item")
	}
	thread, err := store.sessionThread(revised)
	if err != nil || thread != "goal-thread" {
		t.Fatal("goal thread was not resumed", thread, err)
	}
	if _, err := store.decideWorkItem(item.ID, WorkItemDecision{ReviewID: review.ID}, true); err == nil {
		t.Fatal("stale summary approved")
	}
	_, updated := prepareSummary(t, store, revised)
	next, err := store.decideWorkItem(item.ID, WorkItemDecision{ReviewID: updated.Reviews[len(updated.Reviews)-1].ID}, true)
	if err != nil || next == nil {
		t.Fatal(next, err)
	}
	second, _ := store.assignedWorkItem(next.ID)
	if second.Key != "second" {
		t.Fatal("wrong next item", second)
	}
	thread, _ = store.sessionThread(*next)
	if thread != "goal-thread" {
		t.Fatal("next item did not retain goal thread")
	}
	workflow, _ := store.workflow(plan.ID)
	if workflow.Items[0].Status != "approved" || workflow.Items[1].Status != "queued" {
		t.Fatal(workflow)
	}
	nextGoal := startClaimed(t, store)
	_, secondReview := prepareSummary(t, store, nextGoal)
	last, err := store.decideWorkItem(second.ID, WorkItemDecision{ReviewID: secondReview.Reviews[0].ID}, true)
	if err != nil || last != nil {
		t.Fatal("final approval queued extra work", last, err)
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.claimSession(); err != sql.ErrNoRows {
		t.Fatal("approved workflow automatically restarted")
	}
}
func TestSummaryAndCriterionBoundaries(t *testing.T) {
	store := testStore(t)
	plan, _ := submittedPlan(t, store)
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	goal := startClaimed(t, store)
	if err := store.submitAcceptanceCriterion(goal.ID, CriterionInput{Summary: "Missing index", Evidence: []string{"pass"}}); err == nil {
		t.Fatal("missing index accepted")
	}
	if err := store.submitExecutiveSummary(goal.ID, ExecutiveSummaryInput{Summary: "Self approval", Verification: "Claim", Risks: []string{}}); err == nil {
		t.Fatal("goal submitted factory summary")
	}
	submitCriteria(t, store, goal)
	if err := store.requestWorkItemReview(goal.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.claimSession(); err != sql.ErrNoRows {
		t.Fatal("factory reviewer overlapped active goal turn")
	}
	if err := store.finishSession(goal.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	summary := startClaimed(t, store)
	if err := store.submitExecutiveSummary(summary.ID, ExecutiveSummaryInput{Summary: "Summary", Verification: "Checks", Risks: []string{}}); err != nil {
		t.Fatal(err)
	}
	item, _ := store.assignedWorkItem(summary.ID)
	if _, err := store.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID}, true); err == nil {
		t.Fatal("approved summary before factory turn finished")
	}
	index := 0
	if err := store.submitAcceptanceCriterion(summary.ID, CriterionInput{CriterionIndex: &index, Summary: "Cross phase", Evidence: []string{"pass"}}); err == nil {
		t.Fatal("factory modified acceptance submissions")
	}
	if err := store.finishSession(summary.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := store.retryImplementation(goal.StoryID); err == nil {
		t.Fatal("retry bypassed human review")
	}
}
func TestWorkItemHelperProcess(t *testing.T) {
	role := os.Getenv("SAUM_TEST_WORK_ITEM_ROLE")
	if role == "" {
		return
	}
	input, _ := io.ReadAll(os.Stdin)
	args := strings.Join(os.Args, " ")
	if strings.Contains(args, "--ephemeral") {
		fmt.Fprintln(os.Stderr, "goal and factory threads must persist")
		os.Exit(9)
	}
	if expected := os.Getenv("SAUM_TEST_EXPECT_RESUME"); expected != "" && (!strings.Contains(args, "resume") || !strings.Contains(args, expected)) {
		fmt.Fprintln(os.Stderr, "wrong resumed thread", args)
		os.Exit(8)
	}
	thread := "goal-thread"
	if role == "summary" {
		thread = "factory-thread"
	}
	fmt.Printf("{\"type\":\"thread.started\",\"thread_id\":%q}\n", thread)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "work-item-test", Version: "1"}, nil)
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
	if role == "summary" {
		call(`return saum.work_item.submit_executive_summary({summary:"Software is runnable.",verification:"Observed successful go run main.go.",risks:[]});`)
	} else {
		start, end := strings.Index(string(input), "<assigned_work_item>"), strings.Index(string(input), "</assigned_work_item>")
		var item WorkItem
		if start < 0 || end < 0 || json.Unmarshal(input[start+len("<assigned_work_item>"):end], &item) != nil {
			fmt.Fprintln(os.Stderr, "missing one-item assignment")
			os.Exit(7)
		}
		if err := os.WriteFile("main.go", []byte("package main\nimport \"fmt\"\nfunc main(){fmt.Print(\"implemented\")}\n"), 0600); err != nil {
			os.Exit(4)
		}
		output, err := exec.Command("go", "run", "main.go").CombinedOutput()
		if err != nil || string(output) != "implemented" {
			fmt.Fprintln(os.Stderr, string(output), err)
			os.Exit(5)
		}
		for index := range item.AcceptanceCriteria {
			call(fmt.Sprintf(`return saum.work_item.submit_acceptance_criteria({criterionIndex:%d,summary:"Executable works",evidence:["go run main.go printed implemented"]});`, index))
		}
		if role != "skip-request" {
			call(`return saum.work_item.request_review();`)
		}
	}
	_ = session.Close()
	fmt.Println(`{"type":"turn.completed"}`)
	os.Exit(0)
}
func TestGoalAndFactoryAgentMCPRoundTrip(t *testing.T) {
	store := testStore(t)
	plan := twoItemPlan(t, store)
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(store.mcpHandler())
	defer server.Close()
	t.Setenv("SAUM_TEST_MCP_URL", server.URL)
	binary, _ := os.Executable()
	script := filepath.Join(t.TempDir(), "codex")
	_ = os.WriteFile(script, []byte(fmt.Sprintf("#!/bin/sh\nexec %q -test.run=^TestWorkItemHelperProcess$ -- \"$@\"\n", binary)), 0700)
	planner := Planner{Store: store, Root: filepath.Dir(plan.Worktree), Codex: script, MCPURL: server.URL, Timeout: 30 * time.Second}
	for itemIndex := range 2 {
		goal, err := store.claimSession()
		if err != nil {
			t.Fatal(err)
		}
		t.Setenv("SAUM_TEST_WORK_ITEM_ROLE", "goal")
		t.Setenv("SAUM_TEST_EXPECT_RESUME", "")
		if itemIndex > 0 {
			t.Setenv("SAUM_TEST_EXPECT_RESUME", "goal-thread")
		}
		planner.execute(context.Background(), goal)
		finished, _ := store.session(goal.ID)
		if finished.Status != SessionCompleted {
			t.Fatal(finished)
		}
		if err := store.queueApprovedPlans(); err != nil {
			t.Fatal(err)
		}
		summary, err := store.claimSession()
		if err != nil {
			t.Fatal(err)
		}
		t.Setenv("SAUM_TEST_WORK_ITEM_ROLE", "summary")
		t.Setenv("SAUM_TEST_EXPECT_RESUME", "")
		if itemIndex > 0 {
			t.Setenv("SAUM_TEST_EXPECT_RESUME", "factory-thread")
		}
		planner.execute(context.Background(), summary)
		finished, _ = store.session(summary.ID)
		if finished.Status != SessionCompleted {
			t.Fatal(finished)
		}
		item, _ := store.assignedWorkItem(summary.ID)
		if item.Status != "awaiting_review" {
			t.Fatal(item)
		}
		if _, err := store.claimSession(); err != sql.ErrNoRows {
			t.Fatal("next work item leaked before human decision")
		}
		if _, err := store.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[0].ID}, true); err != nil {
			t.Fatal(err)
		}
	}
}
