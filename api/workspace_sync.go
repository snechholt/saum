package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
)

// Only accepted factory commits flow into a story. Checkpoints preserve local
// story work on its own branch; they are not approvals and never advance main.
func (s *Store) syncStoryWorkspace(ctx context.Context, owned OwnedWorkspace) (string, error) {
	s.workspaceMu.Lock()
	defer s.workspaceMu.Unlock()
	main, err := gitOutput(ctx, owned.Repository, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	accepted, err := objectID(main)
	if err != nil {
		return "", err
	}
	note := fmt.Sprintf("Shared factory repository: %s\nStory workspace: %s\nAccepted factory main commit: %s\n", owned.Repository, owned.Worktree, accepted)
	if _, err := gitOutput(ctx, owned.Worktree, "merge-base", "--is-ancestor", accepted, "HEAD"); err == nil {
		return note, nil
	}
	if _, err := gitOutput(ctx, owned.Worktree, "rev-parse", "-q", "--verify", "MERGE_HEAD"); err == nil {
		return note + "A previous merge is unfinished. Resolve its conflicts in this story worktree, commit the merge, and rerun verification before requesting review. Never reset away existing work.\n", nil
	}
	tree, err := captureTree(ctx, owned.Worktree, owned.Repository)
	if err != nil {
		return "", err
	}
	head, err := gitOutput(ctx, owned.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	previous, err := objectID(head)
	if err != nil {
		return "", err
	}
	currentTree, err := gitOutput(ctx, owned.Worktree, "rev-parse", "HEAD^{tree}")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(currentTree)) != string(tree) {
		checkpoint := "refs/saum/checkpoints/" + rand.Text()
		// Keep the original staging tree too, rather than losing partially staged work.
		index, err := gitOutput(ctx, owned.Worktree, "write-tree")
		if err != nil {
			return "", err
		}
		staged, err := objectID(index)
		if err != nil {
			return "", err
		}
		if _, err := gitOutput(ctx, owned.Repository, "update-ref", checkpoint+"-index", staged); err != nil {
			return "", err
		}
		commit, err := snapshotCommit(ctx, owned, tree, previous, "chore: checkpoint story before syncing accepted factory work\n")
		if err != nil {
			return "", err
		}
		if _, err := gitOutput(ctx, owned.Repository, "update-ref", checkpoint, commit); err != nil {
			return "", err
		}
		after, err := captureTree(ctx, owned.Worktree, owned.Repository)
		if err != nil {
			return "", err
		}
		if after != tree {
			return "", fmt.Errorf("story source changed during checkpoint; retry when its writer has stopped")
		}
		branch, err := gitOutput(ctx, owned.Worktree, "symbolic-ref", "HEAD")
		if err != nil {
			return "", err
		}
		if _, err := gitOutput(ctx, owned.Repository, "update-ref", strings.TrimSpace(string(branch)), commit, previous); err != nil {
			return "", err
		}
		if _, err := gitOutput(ctx, owned.Worktree, "read-tree", commit); err != nil {
			return "", err
		}
		note += "Unapproved story work was preserved at " + checkpoint + "; its original staging tree is at " + checkpoint + "-index.\n"
	}
	args := safeCheckoutArgs(ctx, owned.Repository)
	args = append(args, "-c", "user.name=Saum", "-c", "user.email=saum@localhost", "merge", "--ff", "--no-edit", "--no-overwrite-ignore", accepted)
	if _, err := gitOutput(ctx, owned.Worktree, args...); err != nil {
		if _, merging := gitOutput(ctx, owned.Worktree, "rev-parse", "-q", "--verify", "MERGE_HEAD"); merging == nil {
			return note + "Accepted factory changes have merge conflicts with this story. Resolve them here, preserve the intended story changes, commit the merge, and rerun verification. Do not modify factory main.\n", nil
		}
		return "", fmt.Errorf("sync accepted factory changes: %w", err)
	}
	return note + "Accepted factory changes were merged into this story. Reuse this application; do not create a separate replacement.\n", nil
}

func snapshotCommit(ctx context.Context, owned OwnedWorkspace, tree ArtifactTree, parent, message string) (string, error) {
	output, err := gitInput(ctx, owned.Repository, []byte(message), []string{"GIT_AUTHOR_NAME=Saum", "GIT_AUTHOR_EMAIL=saum@localhost", "GIT_COMMITTER_NAME=Saum", "GIT_COMMITTER_EMAIL=saum@localhost"}, "commit-tree", string(tree), "-p", parent)
	if err != nil {
		return "", err
	}
	return objectID(output)
}

// Server-managed checkouts must not run configured hooks, filters, signing, or
// custom merge drivers. Conflicts belong to the story agent, not server hooks.
func safeCheckoutArgs(ctx context.Context, repo string) []string {
	args := []string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "core.attributesFile=/dev/null", "-c", "core.autocrlf=false"}
	configured, _ := gitOutput(ctx, repo, "config", "--null", "--name-only", "--get-regexp", `^(filter\..*\.(process|smudge|clean|required)|merge\..*\.driver)$`)
	for _, key := range strings.Split(string(configured), "\x00") {
		if key == "" {
			continue
		}
		value := ""
		if strings.HasSuffix(key, ".required") || strings.HasSuffix(key, ".driver") {
			value = "false"
		}
		args = append(args, "-c", key+"="+value)
	}
	return args
}
