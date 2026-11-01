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
  const label =
    entry.name === "saum.submit_item"
      ? entry.status === "completed"
        ? "Item saved"
        : "Submit item"
      : entry.name === "saum.complete_plan"
        ? entry.status === "completed"
          ? "Plan submitted"
          : "Complete plan"
        : entry.name === "search"
          ? "Find Saum functions"
          : entry.name === "execute"
            ? "Run Saum functions"
            : entry.name;
  const title = typeof entry.input["title"] === "string" ? entry.input["title"] : "";
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
          <PlanningBadge status={session?.status ?? story.planningStatus} />
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
              label: `Attempt ${sessions.length - index} · ${planningLabel(attempt.status)}`,
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
