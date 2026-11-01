package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Explicit offline operator repair only: old approvals did not capture trees.
// Adopt today's source with a recovery audit event, never fabricate old checks
// or pretend historical reviews were sealed to these bytes.
func (s *Store) recoverApprovedStory(ctx context.Context, planID string) error {
	if err := s.requireOfflineRepair(); err != nil {
		return err
	}
	plan, err := s.session(planID)
	if err != nil {
		return err
	}
	if plan.Kind != PlanningSession || plan.ReviewStatus != ReviewApproved {
		return fmt.Errorf("expected a human-approved plan")
	}
	workflow, err := s.workflow(planID)
	if err != nil {
		return err
	}
	if workflow.Integration != nil && workflow.Integration.Status == "completed" {
		return nil
	}
	if len(workflow.Items) == 0 {
		return fmt.Errorf("story has no approved work")
	}
	reviewID, summaryID := "", ""
	for _, item := range workflow.Items {
		if item.Status != "approved" || len(item.Reviews) == 0 {
			return fmt.Errorf("all story items must already be human-approved")
		}
		review := item.Reviews[len(item.Reviews)-1]
		if review.Decision != ReviewApproved || review.ArtifactTree != "" || review.ApprovedCommit != "" {
			return fmt.Errorf("repair is only for legacy unsealed approvals; use normal review otherwise")
		}
		reviewID, summaryID = review.ID, review.SummarySessionID
	}
	factoryID, _, _, err := s.storyContext(plan.StoryID)
	if err != nil {
		return err
	}
	owned, err := s.ownedWorkspace(ctx, factoryID, plan.Worktree)
	if err != nil {
		return err
	}
	tree, err := captureTree(ctx, owned.Worktree, owned.Repository)
	if err != nil {
		return err
	}
	head, err := gitOutput(ctx, owned.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	parent, err := objectID(head)
	if err != nil {
		return err
	}
	ref := "refs/saum/recoveries/" + planID
	commit := ""
	if output, err := gitOutput(ctx, owned.Repository, "rev-parse", "--verify", ref); err == nil {
		commit, err = objectID(output)
		if err != nil {
			return err
		}
		previous, err := gitOutput(ctx, owned.Repository, "rev-parse", commit+"^{tree}")
		if err != nil || strings.TrimSpace(string(previous)) != string(tree) {
			return fmt.Errorf("legacy source changed since recovery began; inspect the recovery ref")
		}
	} else {
		commit, err = snapshotCommit(ctx, owned, tree, parent, "chore: recover operator-approved legacy story\n\nOriginal legacy approvals were not artifact-bound. Current source adopted during explicit operator-requested workspace repair.\n")
		if err != nil {
			return err
		}
		if _, err := gitOutput(ctx, owned.Repository, "update-ref", ref, commit); err != nil {
			return err
		}
	}
	branch, err := gitOutput(ctx, owned.Worktree, "symbolic-ref", "HEAD")
	if err != nil {
		return err
	}
	if _, err := gitOutput(ctx, owned.Repository, "update-ref", strings.TrimSpace(string(branch)), commit, parent); err != nil {
		return err
	}
	if _, err := gitOutput(ctx, owned.Worktree, "read-tree", commit); err != nil {
		return err
	}
	if err := s.initializeGoalBase(ctx, Session{Kind: ImplementationSession, ParentSessionID: planID}, owned); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM integrations WHERE review_id=?", reviewID).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if err := queueIntegration(tx, planID, reviewID, commit); err != nil {
			return err
		}
	}
	message := fmt.Sprintf("Operator-requested legacy workspace recovery: adopted current source tree %s as commit %s for the already human-approved story. Original unsealed reviews and evidence were preserved without backfilling claims.", tree, commit)
	if _, err := tx.Exec("INSERT INTO session_events(session_id,stream,data) VALUES (?,'workflow',?)", summaryID, []byte(message)); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	if err := s.processIntegrations(); err != nil {
		return err
	}
	result, err := s.goalIntegration(planID)
	if err != nil {
		return err
	}
	if result == nil || result.Status != "completed" {
		return fmt.Errorf("legacy integration needs inspection: %+v", result)
	}
	return nil
}
func (s *Store) requireOfflineRepair() error {
	var active int
	if err := s.db.QueryRow("SELECT count(*) FROM sessions WHERE status IN ('queued','preparing','running')").Scan(&active); err != nil {
		return err
	}
	if active != 0 {
		return fmt.Errorf("stop active work before offline workspace repair")
	}
	return nil
}

// Preserve, rather than delete, a rejected independent bootstrap. This narrow
// repair only accepts an empty original branch with no approved story items.
func (s *Store) archiveRejectedBootstrap(ctx context.Context, planID string) (string, error) {
	if err := s.requireOfflineRepair(); err != nil {
		return "", err
	}
	plan, err := s.session(planID)
	if err != nil {
		return "", err
	}
	workflow, err := s.workflow(planID)
	if err != nil {
		return "", err
	}
	var item *WorkItem
	for i := range workflow.Items {
		candidate := &workflow.Items[i]
		if candidate.Status == "approved" {
			return "", fmt.Errorf("cannot archive a story with accepted work")
		}
		if candidate.Status == "blocked" {
			item = candidate
		}
	}
	if item == nil || len(item.Reviews) == 0 || item.Reviews[len(item.Reviews)-1].Decision != ReviewChangesRequested {
		return "", fmt.Errorf("expected a blocked story with explicit rejected-review feedback")
	}
	factoryID, _, _, err := s.storyContext(plan.StoryID)
	if err != nil {
		return "", err
	}
	owned, err := s.ownedWorkspace(ctx, factoryID, plan.Worktree)
	if err != nil {
		return "", err
	}
	tracked, err := gitOutput(ctx, owned.Worktree, "ls-tree", "-r", "HEAD")
	if err != nil || len(tracked) > 0 {
		return "", fmt.Errorf("only an empty original story branch can be repaired this way")
	}
	index, err := gitOutput(ctx, owned.Worktree, "ls-files")
	if err != nil || len(index) > 0 {
		return "", fmt.Errorf("archive repair requires an empty original index; preserve staged work explicitly first")
	}
	accepted, err := gitOutput(ctx, owned.Repository, "ls-tree", "-r", "HEAD")
	if err != nil || len(accepted) == 0 {
		return "", fmt.Errorf("recover the accepted shared application first")
	}
	tree, err := captureTree(ctx, owned.Worktree, owned.Repository)
	if err != nil {
		return "", err
	}
	head, err := gitOutput(ctx, owned.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	parent, err := objectID(head)
	if err != nil {
		return "", err
	}
	commit, err := snapshotCommit(ctx, owned, tree, parent, "chore: preserve rejected bootstrap before operator workspace repair\n")
	if err != nil {
		return "", err
	}
	recoveryRef := "refs/saum/rejected-bootstraps/" + rand.Text()
	if _, err := gitOutput(ctx, owned.Repository, "update-ref", recoveryRef, commit); err != nil {
		return "", err
	}
	factory, err := os.OpenRoot(owned.Root)
	if err != nil {
		return "", err
	}
	defer factory.Close()
	relative := filepath.Join(".saum-repairs", rand.Text())
	if err := factory.MkdirAll(relative, 0700); err != nil {
		return "", err
	}
	backup, err := canonicalDirectory(filepath.Join(owned.Root, relative))
	if err != nil || !insideDirectory(owned.Root, backup) {
		return "", fmt.Errorf("repair archive must stay inside its factory")
	}
	entries, err := os.ReadDir(owned.Worktree)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		if err := os.Rename(filepath.Join(owned.Worktree, entry.Name()), filepath.Join(backup, entry.Name())); err != nil {
			return backup, err
		}
	}
	if _, err := s.syncStoryWorkspace(ctx, owned); err != nil {
		return backup, err
	}
	note := fmt.Sprintf("Operator-requested workspace repair: the accepted shared application is now available at %s. The rejected independent bootstrap was preserved at %s and %s. Continue the shared application and address the existing review feedback; no item was approved or resumed by this repair.", owned.Worktree, backup, recoveryRef)
	tx, err := s.db.Begin()
	if err != nil {
		return backup, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("INSERT INTO work_item_messages(work_item_id,body) VALUES (?,?)", item.ID, note); err != nil {
		return backup, err
	}
	if _, err := tx.Exec("UPDATE work_items SET blocker='Workspace repaired; awaiting explicit retry to address review feedback.' WHERE id=? AND status='blocked'", item.ID); err != nil {
		return backup, err
	}
	if _, err := tx.Exec("INSERT INTO session_events(session_id,stream,data) VALUES (?,'workflow',?)", item.LastGoalSessionID, []byte(note)); err != nil {
		return backup, err
	}
	return backup, tx.Commit()
}
