import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { Button } from "@cloudflare/kumo/components/button";
import { LayerCard } from "@cloudflare/kumo/components/layer-card";
import { ArrowClockwise, CaretRight } from "@phosphor-icons/react";
import { AgentTranscript } from "./AgentTranscript";
import { getStory, type StoryDetail } from "./story-api";
import {
  cancelSession,
  isActive,
  retryPlanning,
  retryImplementation,
  type Session,
} from "./sessions";
import { usePlanningSession } from "./usePlanningSession";
import { PlanReview } from "./PlanReview";
import { ImplementationProgress } from "./ImplementationProgress";

function StoryView({ detail, sessionId }: { detail: StoryDetail; sessionId: string }) {
  const navigate = useNavigate();
  const session = usePlanningSession(sessionId);
  const [busy, setBusy] = useState(false),
    [error, setError] = useState("");
  const matches = !session.detail || session.detail.session.storyId === detail.story.id;
  const current = matches ? (session.detail?.session ?? null) : null;
  const items = matches ? (session.detail?.items ?? []) : [];
  function selectSession(id: string) {
    void navigate(
      `/stories/${encodeURIComponent(detail.story.id)}/sessions/${encodeURIComponent(id)}`,
    );
  }
  async function retry() {
    setBusy(true);
    setError("");
    try {
      const attempt = await (current?.kind !== "planning" ? retryImplementation : retryPlanning)(
        detail.story.id,
      );
      selectSession(attempt.id);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not start session.");
    } finally {
      setBusy(false);
    }
  }
  async function stop() {
    if (!current) return;
    setBusy(true);
    setError("");
    try {
      await cancelSession(current.id);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not stop session.");
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="story-layout">
      <section className="story-workspace" aria-label="Story workspace">
        <nav className="story-breadcrumb" aria-label="Breadcrumb">
          <Link to="/">Factories</Link>
          <CaretRight size={14} />
          <Link to={`/factories/${encodeURIComponent(detail.factory.id)}`}>
            {detail.factory.name}
          </Link>
        </nav>
        <header className="page-title">
          <h1>{detail.story.title}</h1>
          <Button
            variant="secondary"
            shape="square"
            aria-label={current?.kind !== "planning" ? "Retry work item" : "Plan story"}
            title={current?.kind !== "planning" ? "Retry work item" : "Plan again"}
            loading={busy}
            disabled={
              !detail.story.goal ||
              !current ||
              (current.kind !== "planning" && detail.story.goalStatus !== "failed") ||
              (current.kind === "planning" &&
                (current.reviewStatus !== "pending" ||
                  !!session.detail?.comments.some((comment) => !comment.batchId))) ||
              isActive(detail.story.planningStatus) ||
              busy
            }
            onClick={() => {
              void retry();
            }}
          >
            <ArrowClockwise size={16} />
          </Button>
        </header>
        {error && (
          <p className="error-message" role="alert">
            {error}
          </p>
        )}
        <LayerCard className="story-goal-surface">
          <h2 className="section-title">Goal</h2>
          <p>{detail.story.goal || "—"}</p>
        </LayerCard>
        {matches && current && session.detail?.workflow && (
          <ImplementationProgress
            session={current}
            workflow={session.detail.workflow}
            refresh={session.refresh}
            onSession={selectSession}
          />
        )}
        <PlanReview
          session={matches ? (session.detail?.planSession ?? null) : null}
          items={items}
          comments={matches ? (session.detail?.comments ?? []) : []}
          latest={detail.story.sessionId === sessionId}
          refresh={session.refresh}
          onRevision={(revision: Session) => selectSession(revision.id)}
        />
        {current?.worktree && (
          <details className="story-artifacts">
            <summary>Workspace</summary>
            <code>{current.worktree}</code>
          </details>
        )}
      </section>
      <AgentTranscript
        story={detail.story}
        session={current}
        sessionId={sessionId}
        sessions={detail.sessions}
        events={matches ? session.events : []}
        feedback={matches ? (session.detail?.receivedFeedback ?? []) : []}
        error={matches ? session.error : "Session does not belong to this story."}
        stopping={busy}
        onSelect={selectSession}
        onStop={() => {
          void stop();
        }}
      />
    </div>
  );
}
export default function StoryPage({ id }: { id: string }) {
  const { sessionId = "" } = useParams();
  const [detail, setDetail] = useState<StoryDetail | null>(null),
    [error, setError] = useState("");
  useEffect(() => {
    let ignore = false;
    let timer: ReturnType<typeof setTimeout>;
    async function poll() {
      try {
        const result = await getStory(id);
        if (!ignore) {
          setDetail(result);
          setError("");
        }
      } catch (reason) {
        if (!ignore)
          setError(reason instanceof Error ? reason.message : "Could not load user story.");
      }
      if (!ignore)
        timer = setTimeout(() => {
          void poll();
        }, 1200);
    }
    void poll();
    return () => {
      ignore = true;
      clearTimeout(timer);
    };
  }, [id]);
  if (!detail)
    return (
      <div className="story-loading">
        {error ? (
          <p className="error-message" role="alert">
            {error}
          </p>
        ) : (
          <output>Loading…</output>
        )}
      </div>
    );
  const activeSession = sessionId || detail.story.sessionId;
  return (
    <>
      {error && (
        <p className="error-message" role="alert">
          {error}
        </p>
      )}
      <StoryView key={activeSession || id} detail={detail} sessionId={activeSession} />
    </>
  );
}
