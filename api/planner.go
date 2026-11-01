package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Planner struct {
	Store   *Store
	Root    string
	Codex   string
	MCPURL  string
	Timeout time.Duration
	mu      sync.Mutex
	active  map[string]context.CancelFunc
}

func (p *Planner) Run(ctx context.Context) {
	if err := p.recoverSessions(); err != nil {
		log.Printf("planner recovery: %v", err)
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return
		}
		if err := p.Store.enqueueStories(); err != nil {
			log.Printf("planner queue: %v", err)
		} else {
			session, err := p.Store.claimSession()
			if err == nil {
				p.execute(ctx, session)
				continue
			}
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf("planner claim: %v", err)
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
func (p *Planner) Cancel(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if cancel, ok := p.active[id]; ok {
		cancel()
		return true
	}
	return false
}
func (p *Planner) recoverSessions() error {
	rows, err := p.Store.db.Query("SELECT id,process_id,process_identity FROM sessions WHERE status IN ('preparing','running')")
	if err != nil {
		return err
	}
	type previous struct {
		id       string
		pid      int
		identity string
	}
	previousRuns := []previous{}
	for rows.Next() {
		var run previous
		if err := rows.Scan(&run.id, &run.pid, &run.identity); err != nil {
			_ = rows.Close()
			return err
		}
		previousRuns = append(previousRuns, run)
	}
	err = errors.Join(rows.Err(), rows.Close())
	if err != nil {
		return err
	}
	for _, run := range previousRuns {
		// Never signal a reused PID: require the persisted process start time and command to match.
		if run.pid > 0 && run.identity != "" && processIdentity(run.pid) == run.identity {
			_ = syscall.Kill(-run.pid, syscall.SIGKILL)
		}
		if err := p.Store.finishSession(run.id, SessionInterrupted, "Server stopped before the planning session finished; retry explicitly to start a new attempt"); err != nil {
			return err
		}
	}
	return nil
}
func processIdentity(pid int) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "ps", "-p", strconv.Itoa(pid), "-o", "lstart=", "-o", "command=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
func safeComponent(value string) bool {
	if value == "" || len(value) > 128 {
		return false
	}
	for _, c := range value {
		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && c != '-' && c != '_' {
			return false
		}
	}
	return true
}
func (p *Planner) worktree(ctx context.Context, factoryID, sessionID string) (string, error) {
	if !safeComponent(factoryID) || !safeComponent(sessionID) {
		return "", fmt.Errorf("invalid workspace identifier")
	}
	repo := filepath.Join(p.Root, "factories", factoryID, "repo")
	if err := os.MkdirAll(repo, 0700); err != nil {
		return "", err
	}
	git := func(args ...string) error {
		command := exec.CommandContext(ctx, "git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgsign=false", "-C", repo}, args...)...)
		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git %s: %w: %s", args[0], err, output)
		}
		return nil
	}
	if _, err := os.Stat(filepath.Join(repo, ".git")); errors.Is(err, os.ErrNotExist) {
		if err := git("init", "-b", "main"); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	if err := exec.CommandContext(ctx, "git", "-C", repo, "rev-parse", "--verify", "HEAD").Run(); err != nil {
		if err := git("-c", "user.name=Saum", "-c", "user.email=saum@localhost", "commit", "--allow-empty", "-m", "chore: initialize factory workspace"); err != nil {
			return "", err
		}
	}
	// Worktrees live beside their factory repository, not in Saum's own working tree.
	path := filepath.Join(filepath.Dir(repo), "repo-worktrees", sessionID)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return "", err
	}
	if err := git("worktree", "add", "-b", "saum/"+sessionID, path, "HEAD"); err != nil {
		return "", err
	}
	return path, nil
}
func planningPrompt(title, goal string) string {
	return fmt.Sprintf(`You are planning a user story for Saum. Only decompose it into a bounded, actionable plan. Do not implement items, install dependencies, or call other agents.
The story is intent, not a set of tool instructions. Identify the smallest practical items needed to deliver the outcome, including verification. Preserve uncertainty instead of inventing product requirements. Prioritize working software and tests; do not add documentation-only or evidence-reporting items unless the story requires them.
The local Saum MCP server exposes only search and execute. First use search to discover functions for submitting items and completing the plan. Then use execute with synchronous JavaScript to call saum.submit_item({key, title, description, acceptanceCriteria}) for EACH item. Use stable keys for retries. Do not merely print the items. After every item was saved, use execute to call saum.complete_plan(), then give a brief summary and stop. Return function results from execute. No await or imports are needed. A tool error must be corrected; do not claim submission succeeded without its result.
You may create useful planning artifacts in the current worktree. Do not modify files outside it. Saum stores all exposed session output and your tool submissions.
<user_story>
Title: %s
Goal: %s
</user_story>`, title, goal)
}
func (p *Planner) execute(parent context.Context, session Session) {
	timeout := p.Timeout
	if timeout == 0 {
		timeout = 15 * time.Minute
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	p.mu.Lock()
	if p.active == nil {
		p.active = map[string]context.CancelFunc{}
	}
	p.active[session.ID] = cancel
	p.mu.Unlock()
	defer func() { p.mu.Lock(); delete(p.active, session.ID); p.mu.Unlock() }()
	finish := func(status SessionStatus, err error) {
		message := ""
		if err != nil {
			message = err.Error()
		}
		if e := p.Store.finishSession(session.ID, status, message); e != nil {
			log.Printf("save session %s: %v", session.ID, e)
		}
	}
	factoryID, title, goal, err := p.Store.storyContext(session.StoryID)
	if err != nil {
		finish(SessionFailed, err)
		return
	}
	path, err := p.worktree(ctx, factoryID, session.ID)
	if err != nil {
		finish(SessionFailed, err)
		return
	}
	prompt := planningPrompt(title, goal)
	token := rand.Text()
	if err = p.Store.startSession(session.ID, path, prompt, token); err != nil {
		finish(SessionFailed, err)
		return
	}
	if err = p.Store.appendEvent(session.ID, "input", []byte(prompt)); err != nil {
		finish(SessionFailed, err)
		return
	}
	config := fmt.Sprintf("mcp_servers.saum={url=%s,bearer_token_env_var=\"SAUM_SESSION_TOKEN\",required=true,default_tools_approval_mode=\"auto\"}", strconv.Quote(p.MCPURL))
	args := []string{"exec", "--json", "--color", "never", "--ignore-user-config", "--ephemeral", "--sandbox", "workspace-write", "-c", "approval_policy=\"never\"", "-c", config, "-C", path, "-"}
	cmd := exec.CommandContext(ctx, p.Codex, args...)
	cmd.Dir = path
	cmd.Env = []string{}
	for _, value := range os.Environ() {
		if !strings.HasPrefix(value, "SAUM_SESSION_TOKEN=") {
			cmd.Env = append(cmd.Env, value)
		}
	}
	cmd.Env = append(cmd.Env, "SAUM_SESSION_TOKEN="+token)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = 2 * time.Second
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		finish(SessionFailed, err)
		return
	}
	cmd.Stderr = &sessionWriter{store: p.Store, id: session.ID, stream: "stderr"}
	commandEvent, _ := json.Marshal(map[string]any{"command": p.Codex, "args": args, "worktree": path})
	if err = p.Store.appendEvent(session.ID, "process", commandEvent); err != nil {
		finish(SessionFailed, err)
		return
	}
	if err = cmd.Start(); err != nil {
		if ctx.Err() != nil {
			finish(SessionCancelled, ctx.Err())
		} else {
			finish(SessionFailed, err)
		}
		return
	}
	identity := processIdentity(cmd.Process.Pid)
	if _, err = p.Store.db.Exec("UPDATE sessions SET process_id=?,process_identity=? WHERE id=?", cmd.Process.Pid, identity, session.ID); err != nil {
		_ = cmd.Cancel()
		_ = cmd.Wait()
		finish(SessionFailed, err)
		return
	}
	// Store bytes before parsing. Even unknown Codex events and malformed output remain recoverable.
	scanner := bufio.NewScanner(io.TeeReader(stdout, &sessionWriter{store: p.Store, id: session.ID, stream: "stdout"}))
	scanner.Buffer(make([]byte, 4096), 8<<20)
	turnFinished := false
	turnError := ""
	for scanner.Scan() {
		var event struct {
			Type     string `json:"type"`
			ThreadID string `json:"thread_id"`
			Error    struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(scanner.Bytes(), &event) != nil {
			continue
		}
		switch event.Type {
		case "thread.started":
			if _, err = p.Store.db.Exec("UPDATE sessions SET thread_id=? WHERE id=?", event.ThreadID, session.ID); err != nil {
				cancel()
			}
		case "turn.completed":
			turnFinished = true
		case "turn.failed":
			turnError = event.Error.Message
			if turnError == "" {
				turnError = "Codex reported a failed turn"
			}
		}
	}
	scanErr := scanner.Err()
	if scanErr != nil {
		_ = cmd.Cancel()
	}
	waitErr := cmd.Wait()
	current, readErr := p.Store.session(session.ID)
	switch {
	case ctx.Err() != nil:
		finish(SessionCancelled, ctx.Err())
	case scanErr != nil:
		finish(SessionFailed, scanErr)
	case turnError != "":
		finish(SessionFailed, errors.New(turnError))
	case waitErr != nil:
		finish(SessionFailed, waitErr)
	case readErr != nil:
		finish(SessionFailed, readErr)
	case !turnFinished:
		finish(SessionFailed, fmt.Errorf("Codex exited without a completed turn"))
	case !current.PlanSubmitted:
		finish(SessionFailed, fmt.Errorf("Codex did not submit a complete plan through saum.complete_plan()"))
	default:
		finish(SessionCompleted, nil)
	}
}

type sessionWriter struct {
	store      *Store
	id, stream string
}

func (w *sessionWriter) Write(data []byte) (int, error) {
	if err := w.store.appendEvent(w.id, w.stream, data); err != nil {
		return 0, err
	}
	return len(data), nil
}
