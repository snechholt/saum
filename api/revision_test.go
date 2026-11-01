package main

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFeedbackRevisionRunsInTheExistingWorktree(t *testing.T) {
	store := testStore(t)
	source := testQueuedSession(t, store)
	server := httptest.NewServer(store.mcpHandler())
	defer server.Close()
	t.Setenv("SAUM_TEST_MCP_URL", server.URL)
	planner := Planner{Store: store, Root: t.TempDir(), Codex: fakeCodex(t, "full"), MCPURL: server.URL, Timeout: 20 * time.Second}
	planner.execute(context.Background(), source)
	original, err := store.session(source.ID)
	if err != nil || original.Status != SessionCompleted {
		t.Fatalf("initial plan failed: %+v %v", original, err)
	}
	items, err := store.items(source.ID)
	if err != nil {
		t.Fatal(err)
	}
	comment, err := store.addPlanComment(source.ID, CommentInput{ItemID: items[0].ID, Field: "title", CriterionIndex: -1, Start: 0, End: 9, Quote: "Implement", Body: "Use delivery-oriented wording and make secure sign-in explicit."})
	if err != nil {
		t.Fatal(err)
	}
	revision, err := store.sendPlanFeedback(source.ID)
	if err != nil {
		t.Fatal(err)
	}
	claimed, err := store.claimSession()
	if err != nil || claimed.ID != revision.ID {
		t.Fatal("revision was not claimed", err)
	}
	t.Setenv("SAUM_TEST_FEEDBACK_ID", comment.ID)
	planner.Codex = fakeCodex(t, "revision")
	planner.execute(context.Background(), claimed)
	revised, err := store.session(revision.ID)
	if err != nil || revised.Status != SessionCompleted {
		events, _ := store.events(revision.ID, 0)
		t.Fatalf("revision failed: %+v %v %v", revised, err, events)
	}
	if revised.Worktree != original.Worktree || revised.ReviewStatus != ReviewPending || !strings.Contains(revised.Prompt, comment.Body) {
		t.Fatalf("revision lost context: %+v", revised)
	}
	updated, err := store.items(revision.ID)
	if err != nil || updated[0].Title != "Deliver secure sign-in" {
		t.Fatal("agent did not update copied plan", err)
	}
	previous, err := store.items(source.ID)
	if err != nil || previous[0].Title != "Implement sign-in" {
		t.Fatal("earlier plan was mutated", err)
	}
}
