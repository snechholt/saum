package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

type ArtifactTree string
type snapshotFile struct{ path, mode, object string }

func gitInput(ctx context.Context, directory string, input []byte, env []string, args ...string) ([]byte, error) {
	pinned, err := pinDirectory(directory)
	if err != nil {
		return nil, err
	}
	defer pinned.Close()
	command, err := directoryCommand(ctx, pinned, "git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgsign=false", "-c", "core.fsmonitor=false", "-C", "."}, args...)...)
	if err != nil {
		return nil, err
	}
	command.Env = append(cleanGitEnvironment(), env...)
	command.Stdin = bytes.NewReader(input)
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, output)
	}
	return output, nil
}
func objectID(output []byte) (string, error) {
	id := strings.TrimSpace(string(output))
	if len(id) != 40 && len(id) != 64 {
		return "", fmt.Errorf("invalid Git object identity")
	}
	for _, c := range id {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return "", fmt.Errorf("invalid Git object identity")
		}
	}
	return id, nil
}
func quoteGitPath(path string) string {
	var value strings.Builder
	value.WriteByte('"')
	for _, b := range []byte(path) {
		if b < ' ' || b > 126 || b == '"' || b == '\\' {
			fmt.Fprintf(&value, "\\%03o", b)
		} else {
			value.WriteByte(b)
		}
	}
	value.WriteByte('"')
	return value.String()
}

// Git hashes captured raw bytes in one batch, with collision detection and no
// clean/smudge filters. Neither the goal's staging area nor its refs are changed.
func captureTree(ctx context.Context, worktree, repository string) (ArtifactTree, error) {
	if value, err := gitOutput(ctx, worktree, "config", "--bool", "core.sparseCheckout"); err == nil && strings.TrimSpace(string(value)) == "true" {
		return "", fmt.Errorf("sparse checkouts must be expanded before review")
	}
	paths, err := gitOutput(ctx, worktree, "ls-files", "--cached", "--others", "--exclude-standard", "-z")
	if err != nil {
		return "", err
	}
	root, err := os.OpenRoot(worktree)
	if err != nil {
		return "", err
	}
	defer root.Close()
	staging, err := os.MkdirTemp(filepath.Join(repository, ".git"), "saum-blobs-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(staging)
	files := []snapshotFile{}
	unique := map[string]bool{}
	var batch strings.Builder
	var total int64
	for _, path := range strings.Split(string(paths), "\x00") {
		if path == "" || unique[path] {
			continue
		}
		unique[path] = true
		info, err := root.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}
		if info.Size() > 64<<20 {
			return "", fmt.Errorf("source file exceeds snapshot limit: %s", path)
		}
		total += info.Size()
		if total > 256<<20 {
			return "", fmt.Errorf("source snapshot exceeds 256 MiB")
		}
		var data []byte
		mode := "100644"
		switch {
		case info.Mode().IsRegular():
			file, openErr := root.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
			if openErr != nil {
				return "", openErr
			}
			opened, statErr := file.Stat()
			if statErr != nil || !opened.Mode().IsRegular() {
				_ = file.Close()
				return "", fmt.Errorf("source file changed kind while capturing it")
			}
			data, err = io.ReadAll(io.LimitReader(file, (64<<20)+1))
			_ = file.Close()
			if len(data) > 64<<20 {
				return "", fmt.Errorf("source file grew beyond snapshot limit")
			}
			if info.Mode()&0111 != 0 {
				mode = "100755"
			}
		case info.Mode()&os.ModeSymlink != 0:
			mode = "120000"
			var target string
			target, err = root.Readlink(path)
			if err == nil {
				_, err = root.Stat(path)
			}
			data = []byte(target)
		default:
			return "", fmt.Errorf("unsupported source entry %s; submodules and special files require an explicit full source checkout", path)
		}
		if err != nil {
			return "", fmt.Errorf("source %s: %w", path, err)
		}
		blob := filepath.Join(staging, strconv.Itoa(len(files)))
		if err := os.WriteFile(blob, data, 0600); err != nil {
			return "", err
		}
		batch.WriteString(quoteGitPath(blob))
		batch.WriteByte('\n')
		files = append(files, snapshotFile{path: path, mode: mode})
	}
	if len(files) > 0 {
		output, err := gitInput(ctx, repository, []byte(batch.String()), nil, "hash-object", "-w", "--no-filters", "--stdin-paths")
		if err != nil {
			return "", err
		}
		ids := strings.Fields(string(output))
		if len(ids) != len(files) {
			return "", fmt.Errorf("incomplete snapshot blob capture")
		}
		for index := range files {
			id, err := objectID([]byte(ids[index]))
			if err != nil {
				return "", err
			}
			files[index].object = id
		}
	}
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	index, err := os.CreateTemp(filepath.Join(repository, ".git"), "saum-index-")
	if err != nil {
		return "", err
	}
	name := index.Name()
	_ = index.Close()
	_ = os.Remove(name)
	defer os.Remove(name)
	env := []string{"GIT_INDEX_FILE=" + name}
	if _, err := gitInput(ctx, repository, nil, env, "read-tree", "--empty"); err != nil {
		return "", err
	}
	var entries bytes.Buffer
	for _, file := range files {
		fmt.Fprintf(&entries, "%s %s\t%s\x00", file.mode, file.object, file.path)
	}
	if _, err := gitInput(ctx, repository, entries.Bytes(), env, "update-index", "-z", "--index-info"); err != nil {
		return "", err
	}
	output, err := gitInput(ctx, repository, nil, env, "write-tree")
	if err != nil {
		return "", err
	}
	id, err := objectID(output)
	if err != nil {
		return "", err
	}
	if _, err := gitOutput(ctx, repository, "update-ref", "refs/saum/snapshots/"+id, id); err != nil {
		return "", err
	}
	return ArtifactTree(id), nil
}
func (s *Store) sessionArtifact(ctx context.Context, sessionID string) (OwnedWorkspace, ArtifactTree, error) {
	session, err := s.session(sessionID)
	if err != nil {
		return OwnedWorkspace{}, "", err
	}
	factoryID, _, _, err := s.storyContext(session.StoryID)
	if err != nil {
		return OwnedWorkspace{}, "", err
	}
	plan, err := s.session(session.ParentSessionID)
	if err != nil {
		return OwnedWorkspace{}, "", err
	}
	owned, err := s.ownedWorkspace(ctx, factoryID, plan.Worktree)
	if err != nil {
		return owned, "", err
	}
	if owned.Worktree == owned.Repository {
		return owned, "", fmt.Errorf("goal artifact cannot be the factory main checkout")
	}
	tree, err := captureTree(ctx, owned.Worktree, owned.Repository)
	return owned, tree, err
}
func treeFiles(ctx context.Context, repository string, tree ArtifactTree) ([]snapshotFile, error) {
	output, err := gitOutput(ctx, repository, "ls-tree", "-rz", string(tree))
	if err != nil {
		return nil, err
	}
	files := []snapshotFile{}
	for _, entry := range strings.Split(string(output), "\x00") {
		if entry == "" {
			continue
		}
		metadata, path, ok := strings.Cut(entry, "\t")
		fields := strings.Fields(metadata)
		if !ok || len(fields) != 3 {
			return nil, fmt.Errorf("invalid snapshot entry")
		}
		files = append(files, snapshotFile{path: path, mode: fields[0], object: fields[2]})
	}
	return files, nil
}
func materializeTree(ctx context.Context, repository string, tree ArtifactTree, directory string) error {
	if err := os.Mkdir(directory, 0700); err != nil {
		return err
	}
	files, err := treeFiles(ctx, repository, tree)
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		return err
	}
	defer root.Close()
	pinned, err := pinDirectory(repository)
	if err != nil {
		return err
	}
	defer pinned.Close()
	command, err := directoryCommand(ctx, pinned, "git", "-c", "core.hooksPath=/dev/null", "cat-file", "--batch")
	if err != nil {
		return err
	}
	command.Env = cleanGitEnvironment()
	var requests strings.Builder
	for _, file := range files {
		requests.WriteString(file.object)
		requests.WriteByte('\n')
	}
	command.Stdin = strings.NewReader(requests.String())
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	var diagnostics bytes.Buffer
	command.Stderr = &diagnostics
	if err := command.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	writeErr := func() error {
		for _, file := range files {
			header, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			fields := strings.Fields(header)
			if len(fields) != 3 || fields[0] != file.object || fields[1] != "blob" {
				return fmt.Errorf("invalid snapshot blob")
			}
			size, err := strconv.ParseInt(fields[2], 10, 64)
			if err != nil || size < 0 || size > 64<<20 {
				return fmt.Errorf("invalid snapshot blob size")
			}
			if err := root.MkdirAll(filepath.Dir(file.path), 0700); err != nil {
				return err
			}
			data, err := io.ReadAll(io.LimitReader(reader, size))
			if err != nil || int64(len(data)) != size {
				return fmt.Errorf("incomplete snapshot blob")
			}
			if _, err := reader.ReadByte(); err != nil {
				return err
			}
			switch file.mode {
			case "100644", "100755":
				mode := os.FileMode(0600)
				if file.mode == "100755" {
					mode = 0700
				}
				if err := root.WriteFile(file.path, data, mode); err != nil {
					return err
				}
			case "120000":
				if err := root.Symlink(string(data), file.path); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported snapshot mode")
			}
		}
		return nil
	}()
	if writeErr != nil {
		_ = command.Process.Kill()
	}
	waitErr := command.Wait()
	if writeErr != nil {
		return writeErr
	}
	if waitErr != nil {
		return fmt.Errorf("snapshot materialization: %w: %s", waitErr, diagnostics.String())
	}
	if _, err := gitOutput(ctx, directory, "init", "--template=", "-b", "main"); err != nil {
		return err
	}
	env := []string{"GIT_ALTERNATE_OBJECT_DIRECTORIES=" + filepath.Join(repository, ".git", "objects")}
	_, err = gitInput(ctx, directory, nil, env, "read-tree", string(tree))
	return err
}
