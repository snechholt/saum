package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type OwnedWorkspace struct {
	Root       string
	Repository string
	Worktree   string
	identity   os.FileInfo
}

func insideDirectory(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
func canonicalDirectory(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("not a directory: %s", path)
	}
	return filepath.Abs(resolved)
}
func rejectDirectoryLink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("managed directory must not be a symlink: %s", path)
	}
	return nil
}

func (s *Store) ownedWorkspace(ctx context.Context, factoryID, path string) (OwnedWorkspace, error) {
	root, repository, err := s.factoryWorkspace(factoryID)
	if err != nil {
		return OwnedWorkspace{}, err
	}
	if err := rejectDirectoryLink(root); err != nil {
		return OwnedWorkspace{}, err
	}
	canonicalRoot, err := canonicalDirectory(root)
	if err != nil {
		return OwnedWorkspace{}, err
	}
	repo, err := canonicalDirectory(repository)
	if err != nil {
		return OwnedWorkspace{}, err
	}
	if !insideDirectory(canonicalRoot, repo) || repo == canonicalRoot {
		return OwnedWorkspace{}, fmt.Errorf("repository is outside its factory")
	}
	gitDirectory := filepath.Join(repo, ".git")
	if err := rejectDirectoryLink(gitDirectory); err != nil {
		return OwnedWorkspace{}, err
	}
	if info, err := os.Stat(gitDirectory); err != nil || !info.IsDir() {
		return OwnedWorkspace{}, fmt.Errorf("factory repository must own its Git directory")
	}
	worktree, err := canonicalDirectory(path)
	if err != nil {
		return OwnedWorkspace{}, err
	}
	if !insideDirectory(canonicalRoot, worktree) {
		return OwnedWorkspace{}, fmt.Errorf("workspace resolves outside its factory")
	}
	common, err := gitOutput(ctx, worktree, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return OwnedWorkspace{}, err
	}
	commonPath, err := canonicalDirectory(strings.TrimSpace(string(common)))
	if err != nil {
		return OwnedWorkspace{}, err
	}
	if commonPath != gitDirectory {
		return OwnedWorkspace{}, fmt.Errorf("workspace belongs to a different Git repository")
	}
	registered, err := gitOutput(ctx, repo, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return OwnedWorkspace{}, err
	}
	found := false
	for _, field := range strings.Split(string(registered), "\x00") {
		if strings.HasPrefix(field, "worktree ") {
			candidate, err := canonicalDirectory(strings.TrimPrefix(field, "worktree "))
			if err == nil && candidate == worktree {
				found = true
			}
		}
	}
	if !found {
		return OwnedWorkspace{}, fmt.Errorf("workspace is not registered with its factory repository")
	}
	identity, err := os.Stat(worktree)
	if err != nil {
		return OwnedWorkspace{}, err
	}
	return OwnedWorkspace{Root: canonicalRoot, Repository: repo, Worktree: worktree, identity: identity}, nil
}

// Pin chdir to the checked directory inode: replacing its pathname cannot redirect
// the child into a different directory between validation and process startup.
func (w OwnedWorkspace) pin() (*os.File, error) {
	file, err := os.Open(w.Worktree)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil || !os.SameFile(w.identity, info) {
		_ = file.Close()
		return nil, fmt.Errorf("workspace changed after ownership validation")
	}
	return file, nil
}
func pinDirectory(path string) (*os.File, error) {
	before, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	directory, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	opened, err := directory.Stat()
	if err != nil || !opened.IsDir() || !os.SameFile(before, opened) {
		_ = directory.Close()
		return nil, fmt.Errorf("workspace changed while opening it")
	}
	return directory, nil
}
func gitOutput(ctx context.Context, directory string, args ...string) ([]byte, error) {
	pinned, err := pinDirectory(directory)
	if err != nil {
		return nil, err
	}
	defer pinned.Close()
	command, err := directoryCommand(ctx, pinned, "git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgsign=false", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-C", "."}, args...)...)
	if err != nil {
		return nil, err
	}
	command.Env = cleanGitEnvironment()
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, output)
	}
	return output, nil
}
func cleanGitEnvironment() []string {
	result := []string{}
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if !strings.HasPrefix(key, "GIT_") && key != "SAUM_SESSION_TOKEN" {
			result = append(result, entry)
		}
	}
	return append(result, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0")
}
