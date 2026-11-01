package main

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Opt-in because this uses the local Codex login and makes real model requests.
func TestNativeCodexGoalAndFactoryThreads(t *testing.T) {
	if os.Getenv("SAUM_REAL_CODEX") != "1" {
		t.Skip("set SAUM_REAL_CODEX=1 to exercise native Codex thread resumption")
	}
	store := testStore(t)
	factory, err := store.createFactory("Native thread smoke")
	if err != nil {
		t.Fatal(err)
	}
	spec, err := parseStorySpec("Arithmetic helpers", "Create a tiny Go module with Add and Multiply functions and passing automated tests. No dependencies, services, documentation, or reports.")
	if err != nil {
		t.Fatal(err)
	}
	story, err := store.createStory(factory.ID, spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.enqueueStories(); err != nil {
		t.Fatal(err)
	}
	plan, err := store.claimSession()
	if err != nil {
		t.Fatal(err)
	}
	path := testSessionWorkspace(t, store, plan)
	if err := store.startSession(plan.ID, path, "human-approved test plan", "test-token"); err != nil {
		t.Fatal(err)
	}
	for _, input := range []ItemSubmission{
		{Key: "add", Title: "Implement integer addition", Description: "Initialize a tiny Go module and implement Add(a,b int) int with tests, using only the standard library.", AcceptanceCriteria: []string{"Go tests verify positive, negative, and zero integer addition, including Add(2,3)==5."}},
		{Key: "multiply", Title: "Implement integer multiplication", Description: "Extend the same Go module with Multiply(a,b int) int and tests; preserve Add.", AcceptanceCriteria: []string{"Go tests verify positive, negative, and zero multiplication, including Multiply(3,4)==12, and the existing addition tests still pass."}},
	} {
		if _, err := store.submitItem(plan.ID, input); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.completePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(plan.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(store.mcpHandler())
	defer server.Close()
	planner := Planner{Store: store, Codex: "codex", MCPURL: server.URL, Timeout: 4 * time.Minute}
	var goalThread, factoryThread string
	for index := range 2 {
		run, err := store.claimSession()
		if err != nil {
			t.Fatal(err)
		}
		planner.execute(context.Background(), run)
		run, err = store.session(run.ID)
		if err != nil || run.Status != SessionCompleted {
			t.Fatal("native goal turn failed", run, err)
		}
		if index == 0 {
			goalThread = run.ThreadID
		} else if run.ThreadID != goalThread {
			t.Fatal("goal agent did not resume the same native thread")
		}
		if err := store.queueApprovedPlans(); err != nil {
			t.Fatal(err)
		}
		review, err := store.claimSession()
		if err != nil {
			t.Fatal(err)
		}
		planner.execute(context.Background(), review)
		review, err = store.session(review.ID)
		if err != nil || review.Status != SessionCompleted {
			t.Fatal("factory review failed", review, err)
		}
		if index == 0 {
			factoryThread = review.ThreadID
		} else if review.ThreadID != factoryThread {
			t.Fatal("factory agent did not resume the same native thread")
		}
		item, err := store.assignedWorkItem(review.ID)
		if err != nil || item.Status != "awaiting_review" {
			t.Fatal(item, err)
		}
		t.Logf("item %d waits for human approval: goal thread %s; factory thread %s", index+1, goalThread, factoryThread)
		if _, err := store.decideWorkItem(item.ID, WorkItemDecision{ReviewID: item.Reviews[len(item.Reviews)-1].ID}, true); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.queueApprovedPlans(); err != nil {
		t.Fatal(err)
	}
	stories, err := store.listStories(factory.ID)
	if err != nil || len(stories) != 1 || stories[0].ID != story.ID || stories[0].GoalStatus != "completed" {
		t.Fatal(stories, err)
	}
}
