package main

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Integration struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
	Error  string `json:"error"`
}

func (s *Store) initializeGoalBase(ctx context.Context, session Session, owned OwnedWorkspace) error {
	if session.Kind != ImplementationSession {
		return nil
	}
	var base string
	if err := s.db.QueryRow("SELECT base_commit FROM goal_workflows WHERE plan_session_id=?", session.ParentSessionID).Scan(&base); err != nil {
		return err
	}
	output, err := gitOutput(ctx, owned.Repository, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	main, err := objectID(output)
	if err != nil {
		return err
	}
	output, err = gitOutput(ctx, owned.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	head, err := objectID(output)
	if err != nil {
		return err
	}
	output, err = gitOutput(ctx, owned.Repository, "merge-base", main, head)
	if err != nil {
		return err
	}
	base, err = objectID(output)
	if err != nil {
		return err
	}
	ref, err := gitOutput(ctx, owned.Repository, "symbolic-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("factory main checkout must be on a branch")
	}
	var previousBase, previousRef string
	if err := s.db.QueryRow("SELECT base_commit,base_ref FROM goal_workflows WHERE plan_session_id=?", session.ParentSessionID).Scan(&previousBase, &previousRef); err != nil {
		return err
	}
	branch := strings.TrimSpace(string(ref))
	if previousRef != "" && previousRef != branch {
		return fmt.Errorf("factory main branch changed; restore %s before continuing", previousRef)
	}
	if base != main {
		return fmt.Errorf("merge accepted factory main into this story and finish the merge before requesting review")
	}
	_, err = s.db.Exec("UPDATE goal_workflows SET base_commit=?,base_ref=? WHERE plan_session_id=?", base, branch, session.ParentSessionID)
	return err
}
func approvedCommit(ctx context.Context, owned OwnedWorkspace, tree ArtifactTree, reviewID string) (string, error) {
	branch := "refs/heads/saum/" + filepath.Base(owned.Worktree)
	actual, err := gitOutput(ctx, owned.Worktree, "symbolic-ref", "HEAD")
	if err != nil || strings.TrimSpace(string(actual)) != branch {
		return "", fmt.Errorf("goal branch ownership changed before approval")
	}
	output, err := gitOutput(ctx, owned.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	head, err := objectID(output)
	if err != nil {
		return "", err
	}
	reference := "refs/saum/approved/" + reviewID
	commit := ""
	if existing, err := gitOutput(ctx, owned.Repository, "rev-parse", "--verify", reference); err == nil {
		commit, err = objectID(existing)
		if err != nil {
			return "", err
		}
		actual, err := gitOutput(ctx, owned.Repository, "rev-parse", commit+"^{tree}")
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(string(actual)) != string(tree) {
			return "", fmt.Errorf("approval reference no longer matches reviewed source")
		}
	} else {
		output, err := gitInput(ctx, owned.Repository, []byte("feat: accept reviewed Saum work item\n"), []string{"GIT_AUTHOR_NAME=Saum", "GIT_AUTHOR_EMAIL=saum@localhost", "GIT_COMMITTER_NAME=Saum", "GIT_COMMITTER_EMAIL=saum@localhost"}, "commit-tree", string(tree), "-p", head)
		if err != nil {
			return "", err
		}
		commit, err = objectID(output)
		if err != nil {
			return "", err
		}
		if _, err := gitOutput(ctx, owned.Repository, "update-ref", reference, commit, strings.Repeat("0", len(commit))); err != nil {
			return "", err
		}
	}
	if head != commit {
		parent, err := gitOutput(ctx, owned.Repository, "rev-parse", commit+"^")
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(string(parent)) != head {
			return "", fmt.Errorf("goal branch advanced since approval was prepared; request a fresh review")
		}
		if _, err := gitOutput(ctx, owned.Repository, "update-ref", branch, commit, head); err != nil {
			return "", err
		}
	}
	// No checkout/filter execution: these working bytes were already reviewed.
	if _, err := gitOutput(ctx, owned.Worktree, "read-tree", commit); err != nil {
		return "", err
	}
	return commit, nil
}
func queueIntegration(tx *sql.Tx, planID, reviewID, commit string) error {
	var base, ref string
	if err := tx.QueryRow("SELECT base_commit,base_ref FROM goal_workflows WHERE plan_session_id=?", planID).Scan(&base, &ref); err != nil {
		return err
	}
	if base == "" || ref == "" {
		return fmt.Errorf("goal baseline is missing; reverify the final item before integration")
	}
	_, err := tx.Exec("INSERT INTO integrations(review_id,plan_session_id,base_commit,base_ref,commit_id,status) VALUES (?,?,?,?,?,'queued')", reviewID, planID, base, ref, commit)
	return err
}
func (s *Store) goalIntegration(planID string) (*Integration, error) {
	var result Integration
	err := s.db.QueryRow("SELECT status,commit_id,error FROM integrations WHERE plan_session_id=? ORDER BY rowid DESC LIMIT 1", planID).Scan(&result.Status, &result.Commit, &result.Error)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

// Integration runs only after the operator's final decision is committed. A
// dirty/moved main is never force-reset or automatically conflict-merged.
func (s *Store) processIntegrations() error {
	var reviewID, planID, base, ref, commit, status string
	err := s.db.QueryRow("SELECT review_id,plan_session_id,base_commit,base_ref,commit_id,status FROM integrations WHERE status IN ('queued','running') ORDER BY rowid LIMIT 1").Scan(&reviewID, &planID, &base, &ref, &commit, &status)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	if _, err := s.db.Exec("UPDATE integrations SET status='running' WHERE review_id=?", reviewID); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	runErr := s.integrateSnapshot(ctx, planID, base, ref, commit)
	result, message := "completed", ""
	if runErr != nil {
		result, message = "failed", runErr.Error()
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("UPDATE integrations SET status=?,error=? WHERE review_id=?", result, message, reviewID); err != nil {
		return err
	}
	if runErr == nil {
		if _, err := tx.Exec("UPDATE goal_workflows SET integrated_commit=? WHERE plan_session_id=?", commit, planID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec("INSERT INTO session_events(session_id,stream,data) SELECT summary_session_id,'workflow',? FROM work_item_reviews WHERE id=?", []byte(fmt.Sprintf("Integration %s: %s %s", result, commit, message)), reviewID); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *Store) integrateSnapshot(ctx context.Context, planID, base, ref, commit string) error {
	s.workspaceMu.Lock()
	defer s.workspaceMu.Unlock()
	plan, err := s.session(planID)
	if err != nil {
		return err
	}
	factoryID, _, _, err := s.storyContext(plan.StoryID)
	if err != nil {
		return err
	}
	_, repo, err := s.factoryWorkspace(factoryID)
	if err != nil {
		return err
	}
	owned, err := s.ownedWorkspace(ctx, factoryID, repo)
	if err != nil {
		return err
	}
	actualRef, err := gitOutput(ctx, owned.Repository, "symbolic-ref", "HEAD")
	if err != nil || strings.TrimSpace(string(actualRef)) != ref {
		return fmt.Errorf("factory main branch changed; restore the intended branch and reverify")
	}
	head, err := gitOutput(ctx, owned.Repository, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	current, err := objectID(head)
	if err != nil {
		return err
	}
	dirty, err := gitOutput(ctx, owned.Repository, "status", "--porcelain", "-z", "--untracked-files=normal")
	if err != nil {
		return err
	}
	if len(dirty) != 0 {
		return fmt.Errorf("factory main has local changes; integration stopped without overwriting them")
	}
	if current == commit {
		return nil
	} // Recover a crash after Git's atomic ref update.
	if current != base {
		return fmt.Errorf("factory main advanced; rework the final item, rebase its goal branch onto main, and rerun verification")
	}
	if _, err := gitOutput(ctx, owned.Repository, "merge-base", "--is-ancestor", base, commit); err != nil {
		return fmt.Errorf("reviewed commit does not fast-forward the goal baseline")
	}
	// Do not execute repository-configured filter programs with server privileges.
	filters, _ := gitOutput(ctx, owned.Repository, "config", "--null", "--name-only", "--get-regexp", `^filter\..*\.(process|smudge|clean|required)$`)
	args := []string{"-c", "core.attributesFile=/dev/null", "-c", "core.autocrlf=false"}
	for _, key := range strings.Split(string(filters), "\x00") {
		if key == "" {
			continue
		}
		value := ""
		if strings.HasSuffix(key, ".required") {
			value = "false"
		}
		args = append(args, "-c", key+"="+value)
	}
	args = append(args, "merge", "--ff-only", "--no-edit", commit)
	if _, err := gitOutput(ctx, owned.Repository, args...); err != nil {
		return err
	}
	head, err = gitOutput(ctx, owned.Repository, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(head)) != commit {
		return fmt.Errorf("main changed concurrently; integration needs operator inspection")
	}
	return nil
}
