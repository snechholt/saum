package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func submittedPlan(t *testing.T, store *Store) (Session, ActionItem) {
	t.Helper()
	session := testQueuedSession(t, store)
	if err := store.startSession(session.ID, testSessionWorkspace(t, store, session), "prompt", "secret"); err != nil {
		t.Fatal(err)
	}
	item, err := store.submitItem(session.ID, ItemSubmission{Key: "login", Title: "Login 🎯", Description: "Users sign in securely.", AcceptanceCriteria: []string{"Valid credentials open the private area", "Invalid credentials are rejected"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.completePlan(session.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(session.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	session, err = store.session(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	return session, item
}
func titleComment(item ActionItem) CommentInput {
	return CommentInput{ItemID: item.ID, Field: "title", CriterionIndex: -1, Start: 0, End: 5, Quote: "Login", Body: "Make the title describe username/password sign-in."}
}
func TestPlanApprovalIsAnExplicitImmutableBatch(t *testing.T) {
	store := testStore(t)
	session, item := submittedPlan(t, store)
	if err := store.approvePlan(session.ID); err != nil {
		t.Fatal(err)
	}
	approved, err := store.session(session.ID)
	if err != nil || approved.ReviewStatus != ReviewApproved || approved.ReviewedAt == "" {
		t.Fatalf("not approved: %+v %v", approved, err)
	}
	if _, err := store.addPlanComment(session.ID, titleComment(item)); !errors.Is(err, ErrReviewConflict) {
		t.Fatal("annotated approved plan", err)
	}
	if _, err := store.retryStory(session.StoryID); err == nil {
		t.Fatal("replanned approved state", err)
	}
	if _, err := store.updateItem(session.ID, ItemSubmission{Key: "login", Title: "Changed", Description: "Changed", AcceptanceCriteria: []string{"Changed"}}); err != ErrPlanClosed {
		t.Fatal("mutated approved plan", err)
	}
}
func TestPlanCommentsPinExactTextAndBlockApproval(t *testing.T) {
	store := testStore(t)
	session, item := submittedPlan(t, store)
	input := titleComment(item)
	input.Start = 6
	input.End = 8
	input.Quote = "🎯"
	comment, err := store.addPlanComment(session.ID, input)
	if err != nil {
		t.Fatal("UTF-16 emoji selection failed", err)
	}
	if err := store.approvePlan(session.ID); !errors.Is(err, ErrReviewConflict) {
		t.Fatal("approved with pending comments", err)
	}
	for _, bad := range []CommentInput{
		{ItemID: item.ID, Field: "title", CriterionIndex: -1, Start: 6, End: 7, Quote: "🎯", Body: "Wrong boundary"},
		{ItemID: item.ID, Field: "title", CriterionIndex: -1, Start: 0, End: 5, Quote: "Other", Body: "Stale quote"},
		{ItemID: item.ID, Field: "acceptanceCriteria", CriterionIndex: 20, Start: 0, End: 5, Quote: "Valid", Body: "Missing criterion"},
		{ItemID: item.ID, Field: "goal", CriterionIndex: -1, Start: 0, End: 5, Quote: "Login", Body: "Wrong field"},
	} {
		if _, err := store.addPlanComment(session.ID, bad); err == nil {
			t.Fatalf("accepted invalid comment: %+v", bad)
		}
	}
	if err := store.removePlanComment(session.ID, comment.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.approvePlan(session.ID); err != nil {
		t.Fatal("removing last comment should allow approval", err)
	}
}
func TestFeedbackQueuesOneRevisionAndPreservesOriginalPlan(t *testing.T) {
	store := testStore(t)
	source, item := submittedPlan(t, store)
	first, err := store.addPlanComment(source.ID, titleComment(item))
	if err != nil {
		t.Fatal(err)
	}
	criterion := CommentInput{ItemID: item.ID, Field: "acceptanceCriteria", CriterionIndex: 1, Start: 0, End: 7, Quote: "Invalid", Body: "Also cover missing credentials."}
	second, err := store.addPlanComment(source.ID, criterion)
	if err != nil {
		t.Fatal(err)
	}
	revision, err := store.sendPlanFeedback(source.ID)
	if err != nil {
		t.Fatal(err)
	}
	duplicate, err := store.sendPlanFeedback(source.ID)
	if err != nil || duplicate.ID != revision.ID {
		t.Fatalf("batch sent twice: %+v %v", duplicate, err)
	}
	if revision.ParentSessionID != source.ID || revision.FeedbackBatchID == "" {
		t.Fatalf("revision missing provenance: %+v", revision)
	}
	comments, err := store.receivedFeedback(revision.FeedbackBatchID)
	if err != nil || len(comments) != 2 {
		t.Fatalf("feedback not batched: %+v %v", comments, err)
	}
	cloned, err := store.items(revision.ID)
	if err != nil || len(cloned) != 1 || cloned[0].ID == item.ID || cloned[0].Key != item.Key {
		t.Fatalf("plan not copied safely: %+v %v", cloned, err)
	}
	if err := store.approvePlan(source.ID); !errors.Is(err, ErrReviewConflict) {
		t.Fatal("approved superseded plan", err)
	}
	claimed, err := store.claimSession()
	if err != nil || claimed.ID != revision.ID {
		t.Fatal("revision not queued", err)
	}
	if err := store.startSession(revision.ID, t.TempDir(), "revision", "next-secret"); err != nil {
		t.Fatal(err)
	}
	if err := store.completePlan(revision.ID); err == nil {
		t.Fatal("ignored unresolved feedback")
	}
	code := `return saum.update_item({key:"login",title:"Username/password sign-in",description:"Users sign in securely.",acceptanceCriteria:["Valid credentials open the private area","Invalid and missing credentials are rejected"]});`
	if _, err := store.executeCode(context.Background(), revision.ID, code); err != nil {
		t.Fatal(err)
	}
	for _, comment := range []PlanComment{first, second} {
		if err := store.resolveFeedback(revision.ID, FeedbackResolution{CommentID: comment.ID, Resolution: "Updated the item and its acceptance criteria."}); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.resolveFeedback(revision.ID, FeedbackResolution{CommentID: "wrong", Resolution: "Not assigned"}); err == nil {
		t.Fatal("resolved unassigned feedback")
	}
	if err := store.completePlan(revision.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(revision.ID, SessionCompleted, ""); err != nil {
		t.Fatal(err)
	}
	original, err := store.items(source.ID)
	if err != nil || original[0].Title != "Login 🎯" {
		t.Fatal("original plan changed", err)
	}
	updated, err := store.items(revision.ID)
	if err != nil || updated[0].Title != "Username/password sign-in" {
		t.Fatal("revision not updated", err)
	}
	if err := store.approvePlan(revision.ID); err != nil {
		t.Fatal(err)
	}
}
func TestFailedRevisionRetryKeepsItsFeedbackAndBase(t *testing.T) {
	store := testStore(t)
	source, item := submittedPlan(t, store)
	if _, err := store.addPlanComment(source.ID, titleComment(item)); err != nil {
		t.Fatal(err)
	}
	revision, err := store.sendPlanFeedback(source.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.claimSession(); err != nil {
		t.Fatal(err)
	}
	if err := store.finishSession(revision.ID, SessionFailed, "Provider unavailable"); err != nil {
		t.Fatal(err)
	}
	retry, err := store.retryStory(source.StoryID)
	if err != nil {
		t.Fatal(err)
	}
	if retry.ParentSessionID != source.ID || retry.FeedbackBatchID != revision.FeedbackBatchID {
		t.Fatal("retry lost feedback context", retry)
	}
	items, err := store.items(retry.ID)
	if err != nil || len(items) != 1 || items[0].Title != item.Title {
		t.Fatal("retry lost source plan", err)
	}
}
func TestFeedbackNeverCrossesPlanBoundaries(t *testing.T) {
	store := testStore(t)
	first, item := submittedPlan(t, store)
	second, _ := submittedPlan(t, store)
	if _, err := store.addPlanComment(second.ID, titleComment(item)); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal("comment linked to another plan's item", err)
	}
	if _, err := store.sendPlanFeedback(first.ID); !errors.Is(err, ErrReviewConflict) {
		t.Fatal("sent an empty feedback batch", err)
	}
}
