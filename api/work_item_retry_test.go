package main

import "testing"

func TestRetryBlockedWorkItemWithoutReplyPreservesGoalAndApprovalGate(t *testing.T) {
	store := testStore(t)
	plan := twoItemPlan(t, store)
	if err := store.approvePlan(plan.ID); err != nil {
		t.Fatal(err)
	}
	previous := startClaimed(t, store)
	thread, err := store.sessionThread(previous)
	if err != nil {
		t.Fatal(err)
	}
	if thread == "" {
		thread = "existing-goal-thread"
		if err := store.saveSessionThread(previous, thread); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.blockWorkItem(previous.ID, "Startup check could not bind localhost"); err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(previous.ID, SessionFailed, "Blocked"); err != nil {
		t.Fatal(err)
	}
	next, err := store.retryImplementation(previous.StoryID)
	if err != nil {
		t.Fatal(err)
	}
	if next.ID == previous.ID || next.ParentSessionID != plan.ID || next.Kind != ImplementationSession {
		t.Fatal(next)
	}
	resumed, err := store.sessionThread(next)
	if err != nil || resumed != thread {
		t.Fatalf("thread changed: %s %v", resumed, err)
	}
	item, err := store.assignedWorkItem(next.ID)
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != "queued" || len(item.Replies) != 0 || item.CurrentSessionID != next.ID {
		t.Fatal(item)
	}
	workflow, err := store.workflow(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(workflow.Items) != 2 || workflow.Items[1].Status != "pending" {
		t.Fatal(workflow)
	}
	if _, err := store.retryImplementation(previous.StoryID); err == nil {
		t.Fatal("duplicate retry queued active work")
	}
}
