import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router";
import { AgentTranscript } from "./AgentTranscript";
import { getStory, type StoryDetail } from "./story-api";
import { usePlanningSession } from "./usePlanningSession";
import { auditPath, sessionPath } from "./paths";
import type { SessionEvent } from "./sessions";
function rawEvents(events: SessionEvent[]): string {
  const decoders = new Map<string, TextDecoder>();
  return events
    .map((event) => {
      let decoder = decoders.get(event.stream);
      if (!decoder) {
        decoder = new TextDecoder();
        decoders.set(event.stream, decoder);
      }
      return `[${event.id} ${event.stream}]\n${decoder.decode(
        Uint8Array.from(atob(event.data), (character) => character.charCodeAt(0)),
        { stream: true },
      )}`;
    })
    .join("\n");
}
export function AuditTranscript({ id }: { id: string }) {
  const snapshot = usePlanningSession(id),
    storyId = snapshot.detail?.session.storyId ?? "",
    [story, setStory] = useState<StoryDetail | null>(null),
    [error, setError] = useState(""),
    [rawOpen, setRawOpen] = useState(false),
    navigate = useNavigate();
  useEffect(() => {
    if (!storyId) return;
    let ignore = false;
    void getStory(storyId)
      .then((result) => {
        if (!ignore) setStory(result);
      })
      .catch((reason) => {
        if (!ignore) setError(reason instanceof Error ? reason.message : "Could not load story.");
      });
    return () => {
      ignore = true;
    };
  }, [storyId]);
  const current = snapshot.detail?.session;
  return (
    <section className="audit-detail">
      <nav className="story-breadcrumb">
        <Link to="/audit">Audit</Link>
        {story && (
          <Link to={sessionPath(story.factory.id, story.story.id, id)}>Open in factory</Link>
        )}
      </nav>
      {(error || snapshot.error) && (
        <p className="error-message" role="alert">
          {error || snapshot.error}
        </p>
      )}
      {current && story ? (
        <>
          <AgentTranscript
            story={story.story}
            session={current}
            sessionId={id}
            sessions={story.sessions}
            events={snapshot.events}
            feedback={snapshot.detail?.receivedFeedback ?? []}
            error=""
            stopping={false}
            onSelect={(id) => {
              void navigate(auditPath(id));
            }}
          />
          <details className="audit-session-details">
            <summary>Session details</summary>
            <dl>
              <dt>Session</dt>
              <dd>
                <code>{current.id}</code>
              </dd>
              <dt>Codex thread</dt>
              <dd>
                <code>{current.threadId || "—"}</code>
              </dd>
              <dt>Workspace</dt>
              <dd>
                <code>{current.worktree || "—"}</code>
              </dd>
            </dl>
            <details>
              <summary>Prompt</summary>
              <pre>{current.prompt}</pre>
            </details>
            <details onToggle={(event) => setRawOpen(event.currentTarget.open)}>
              <summary>Raw events</summary>
              <pre>{rawOpen ? rawEvents(snapshot.events) : ""}</pre>
            </details>
          </details>
        </>
      ) : (
        !snapshot.error && <output>Loading…</output>
      )}
    </section>
  );
}
