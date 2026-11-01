import { useEffect, useMemo, useRef } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Select } from "@cloudflare/kumo/components/select";
import {
  CheckCircle,
  Robot,
  SpinnerGap,
  Stop,
  Terminal,
  Wrench,
  XCircle,
} from "@phosphor-icons/react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { UserStory } from "./api";
import type { PlanComment } from "./sessions";
import { isActive, planningLabel, type Session, type SessionEvent } from "./sessions";
import { PlanningBadge } from "./PlanningBadge";
import { parseTranscript, type TranscriptEntry } from "./transcript";

const markdownPlugins = [remarkGfm];
function ToolStatus({ status }: { status: string }) {
  return status === "failed" ? (
    <XCircle size={16} className="tool-failed" />
  ) : status === "completed" ? (
    <CheckCircle size={16} className="tool-complete" />
  ) : (
    <SpinnerGap size={16} className="spin" />
  );
}
function ToolEntry({ entry }: { entry: Exclude<TranscriptEntry, { kind: "message" }> }) {
  if (entry.kind === "command")
    return (
      <details className="agent-tool">
        <summary>
          <Terminal size={16} />
          <span>Run command</span>
          <ToolStatus status={entry.status} />
        </summary>
        <pre>{entry.command}</pre>
        {entry.output && <pre>{entry.output}</pre>}
      </details>
    );
  const labels: Record<string, { pending: string; done: string }> = {
    "saum.submit_item": { pending: "Submit item", done: "Item saved" },
    "saum.update_item": { pending: "Update item", done: "Item updated" },
    "saum.remove_item": { pending: "Remove item", done: "Item removed" },
    "saum.resolve_feedback": { pending: "Address feedback", done: "Feedback addressed" },
    "saum.start_item": { pending: "Start item", done: "Item started" },
    "saum.complete_item": { pending: "Complete item", done: "Item completed" },
    "saum.block_item": { pending: "Report blocker", done: "Item blocked" },
    "saum.complete_implementation": {
      pending: "Complete implementation",
      done: "Implementation submitted",
    },
    "saum.complete_plan": { pending: "Complete plan", done: "Plan submitted" },
    search: { pending: "Find Saum functions", done: "Find Saum functions" },
    execute: { pending: "Run Saum functions", done: "Run Saum functions" },
  };
  const names = labels[entry.name],
    label = names ? (entry.status === "completed" ? names.done : names.pending) : entry.name;
  const title =
    typeof entry.input["title"] === "string"
      ? entry.input["title"]
      : typeof entry.input["resolution"] === "string"
        ? entry.input["resolution"]
        : typeof entry.input["summary"] === "string"
          ? entry.input["summary"]
          : typeof entry.input["key"] === "string"
            ? entry.input["key"]
            : "";
  const code =
    typeof entry.input["code"] === "string"
      ? entry.input["code"]
      : JSON.stringify(entry.input, null, 2);
  return (
    <details className="agent-tool">
      <summary>
        <Wrench size={16} />
        <span>
          {label}
          {title && <span className="tool-subtitle">{title}</span>}
        </span>
        <ToolStatus status={entry.status} />
      </summary>
      <pre>{code}</pre>
      {entry.error && <p className="error-message">{entry.error}</p>}
      {entry.result !== undefined && entry.result !== null && (
        <pre>{JSON.stringify(entry.result, null, 2)}</pre>
      )}
    </details>
  );
}
export function AgentTranscript({
  story,
  session,
  sessionId,
  sessions,
  events,
  feedback,
  error,
  stopping,
  onSelect,
  onStop,
}: {
  story: UserStory;
  session: Session | null;
  sessionId: string;
  sessions: Session[];
  events: SessionEvent[];
  feedback: PlanComment[];
  error: string;
  stopping: boolean;
  onSelect: (id: string) => void;
  onStop: () => void;
}) {
  const parsed = useMemo(() => parseTranscript(events), [events]);
  const scroll = useRef<HTMLDivElement>(null),
    following = useRef(true);
  useEffect(() => {
    if (parsed.entries.length && following.current && scroll.current)
      scroll.current.scrollTop = scroll.current.scrollHeight;
  }, [parsed.entries]);
  const status = session?.status ?? story.planningStatus;
  const active = isActive(status);
  return (
    <aside className="agent-panel" aria-label="Agent transcript">
      <header className="agent-panel-header">
        <div>
          <Robot size={20} />
          <strong>Codex</strong>
        </div>
        <div className="session-actions">
          <PlanningBadge
            status={session?.status ?? story.planningStatus}
            kind={session?.kind ?? story.sessionKind}
          />
          {session && active && (
            <Button
              variant="ghost"
              shape="square"
              aria-label="Stop session"
              title="Stop session"
              loading={stopping}
              onClick={onStop}
            >
              <Stop size={16} />
            </Button>
          )}
        </div>
      </header>
      {sessions.length > 0 && (
        <div className="agent-session-picker">
          <Select<string>
            aria-label="Session attempt"
            value={sessionId || null}
            onValueChange={(value) => {
              if (value) onSelect(value);
            }}
            items={sessions.map((attempt, index) => ({
              value: attempt.id,
              label: `${attempt.kind === "implementation" ? "Build" : "Plan"} ${sessions.length - index} · ${planningLabel(attempt.status, attempt.kind)}`,
            }))}
          />
        </div>
      )}
      <div
        className="agent-scroll"
        ref={scroll}
        onScroll={() => {
          const element = scroll.current;
          if (element)
            following.current =
              element.scrollHeight - element.scrollTop - element.clientHeight < 80;
        }}
      >
        <div className="story-chat-message">
          <span className="message-label">User story</span>
          <p>{story.title}</p>
          <div>{story.goal || "No goal specified."}</div>
        </div>
        {feedback.length > 0 && (
          <div className="story-chat-message review-chat-message">
            <span className="message-label">Feedback</span>
            {feedback.map((comment) => (
              <div className="review-chat-comment" key={comment.id}>
                <blockquote>{comment.quote}</blockquote>
                <p>{comment.body}</p>
              </div>
            ))}
          </div>
        )}
        {parsed.entries.map((entry) =>
          entry.kind === "message" ? (
            <article className="assistant-message" key={entry.id}>
              <div className="message-label">
                <Robot size={16} />
                Codex
              </div>
              <div className="chat-markdown">
                <Markdown remarkPlugins={markdownPlugins}>{entry.text}</Markdown>
              </div>
            </article>
          ) : (
            <ToolEntry key={entry.id} entry={entry} />
          ),
        )}
        {active && (
          <div className="agent-working">
            <SpinnerGap size={16} className="spin" />
            {status === "queued"
              ? "Queued"
              : status === "preparing"
                ? "Preparing workspace…"
                : "Working…"}
          </div>
        )}
        {(error || session?.error) && (
          <p className="error-message" role="alert">
            {error || session?.error}
          </p>
        )}
        {!sessionId && !active && <p className="agent-empty">No session yet.</p>}
        {parsed.diagnostics && (
          <details className="agent-diagnostics">
            <summary>Diagnostics</summary>
            <pre>{parsed.diagnostics}</pre>
          </details>
        )}
      </div>
    </aside>
  );
}
