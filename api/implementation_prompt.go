package main

import (
	"encoding/json"
	"fmt"
)

func (s *Store) implementationPrompt(session Session, title, goal string) (string, error) {
	plan, err := s.session(session.ParentSessionID)
	if err != nil {
		return "", err
	}
	if plan.Kind != PlanningSession || plan.ReviewStatus != ReviewApproved || plan.StoryID != session.StoryID {
		return "", fmt.Errorf("implementation does not have an approved plan")
	}
	items, err := s.items(plan.ID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`Implement this APPROVED user-story plan. Approval means build working software now, not another plan or a documentation report.
Title: %s
Goal: %s

You are in the factory's isolated Git worktree. Inspect its files first, reuse the existing stack and conventions, and preserve previous work. If there is no application, build the smallest runnable application that delivers the approved scope. Do not modify files outside this worktree, deploy, access unrelated credentials, or create a planning/evidence report instead of software. Do not change the approved plan.

Use the local MCP search tool to discover implementation progress functions. For each approved item, call saum.start_item({key}) through execute BEFORE beginning it. Implement and run the relevant tests, then call saum.complete_item({key,summary}) with a concise description of the changes and actual verification results. Work through items sequentially in the plan order. Functions use synchronous JavaScript; return results, no await/imports.
If a necessary product decision or external dependency genuinely blocks safe implementation, call saum.block_item({key,summary}) with the exact question or blocker and stop. Do not invent security-sensitive account ownership, recovery channels, or production credentials; do not claim incomplete work is done.
Only after EVERY item is implemented and verified, call saum.complete_implementation({summary}). Include actual test results and the command to run the application in the summary. Keep the application and tests in this worktree. Finish with a brief user-facing summary. Saum persists the transcript, progress, and worktree; successful process exit without completion does NOT count as implementation.

<approved_plan>
%s
</approved_plan>`, title, goal, data), nil
}
