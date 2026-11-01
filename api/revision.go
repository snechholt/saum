package main

import (
	"context"
	"encoding/json"
	"fmt"
)

func (p *Planner) sessionWorkspace(ctx context.Context, session Session, factoryID string) (string, error) {
	if session.ParentSessionID == "" {
		return p.worktree(ctx, factoryID, session.ID)
	}
	parent, err := p.Store.session(session.ParentSessionID)
	if err != nil {
		return "", err
	}
	if parent.StoryID != session.StoryID || parent.Status != SessionCompleted || !parent.PlanSubmitted {
		return "", fmt.Errorf("revision source is not a submitted plan for this story")
	}
	owned, err := p.Store.ownedWorkspace(ctx, factoryID, parent.Worktree)
	if err != nil {
		return "", err
	}
	if session.Kind == SummarySession {
		return owned.Repository, nil
	}
	return owned.Worktree, nil
}
func (s *Store) revisionPrompt(session Session, title, goal string) (string, error) {
	items, err := s.items(session.ParentSessionID)
	if err != nil {
		return "", err
	}
	comments, err := s.receivedFeedback(session.FeedbackBatchID)
	if err != nil {
		return "", err
	}
	if len(comments) == 0 {
		return "", fmt.Errorf("revision has no feedback")
	}
	context, err := json.MarshalIndent(struct {
		Items    []ActionItem  `json:"sourcePlan"`
		Comments []PlanComment `json:"feedback"`
	}{items, comments}, "", "  ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`Revise the existing user-story plan in Saum in response to the human's feedback. Do not implement the story, install dependencies, or create documentation/reporting work that was not requested.
Title: %s
Goal: %s

Your working plan is already a COPY of the previous plan, with the same item keys. The previous plan and comments are immutable. Do not regenerate or duplicate the whole plan. Keep unaffected items; update the necessary ones.
Use the local MCP search tool to discover saum.update_item, saum.remove_item, saum.submit_item, saum.resolve_feedback, and saum.complete_plan. Then use execute with synchronous JavaScript. Return results; no imports or await.
Use saum.update_item with all fields to replace an existing item by key. Use saum.submit_item only for genuinely new keys and saum.remove_item for obsolete ones.
After saving the actual changes for EACH assigned feedback comment, call saum.resolve_feedback({commentId, resolution}) describing how it was addressed. Every comment must have a resolution before saum.complete_plan() can succeed. If a request conflicts with the goal, explain the conflict in its resolution rather than silently changing the goal.
Finally call saum.complete_plan(), briefly summarize the changes, and stop. This only submits a revised plan; the human must still approve it.

<review_context>
%s
</review_context>`, title, goal, context), nil
}
