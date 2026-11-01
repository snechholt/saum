import { useState } from "react";
import { Badge } from "@cloudflare/kumo/components/badge";
import { Button } from "@cloudflare/kumo/components/button";
import { InputArea } from "@cloudflare/kumo/components/input";
import { LayerCard } from "@cloudflare/kumo/components/layer-card";
import { ArrowClockwise, Check, ChatText, Play } from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {
  retryImplementation,
  type GoalWorkflow,
  type Session,
  type WorkItem,
  type WorkItemStatus,
} from "./sessions";
import { decideWorkItem, replyToWorkItem } from "./work-item-api";
const labels: Record<WorkItemStatus, string> = {
  pending: "Pending",
  queued: "Queued",
  working: "Working",
  blocked: "Blocked",
  review_requested: "Review requested",
  summarizing: "Summarizing",
  awaiting_review: "Needs review",
  approved: "Approved",
  failed: "Failed",
};
function ItemBadge({ status }: { status: WorkItemStatus }) {
  return (
    <Badge
      variant={
        status === "approved"
          ? "success"
          : status === "failed" || status === "blocked"
            ? "error"
            : status === "pending"
              ? "neutral"
              : "info"
      }
    >
      {labels[status]}
    </Badge>
  );
}
const columns: Column<WorkItem>[] = [
  { id: "title", header: "Item", value: (item) => item.title },
  {
    id: "status",
    header: "Status",
    value: (item) => item.status,
    render: (item) => <ItemBadge status={item.status} />,
  },
  {
    id: "criteria",
    header: "Criteria",
    value: (item) => item.criteria.length,
    render: (item) => (
      <span>
        {item.criteria.length} / {item.acceptanceCriteria.length}
      </span>
    ),
  },
];
const plugins = [remarkGfm];
export function ImplementationProgress({
  session,
  workflow,
  refresh,
  onSession,
}: {
  session: Session;
  workflow: GoalWorkflow;
  refresh: () => Promise<void>;
  onSession: (id: string) => void;
}) {
  const [feedback, setFeedback] = useState(""),
    [reply, setReply] = useState(""),
    [busy, setBusy] = useState(false),
    [error, setError] = useState("");
  const active =
    workflow.items.find(
      (item) =>
        item.currentSessionId === session.id ||
        item.reviews.some(
          (review) => review.summarySessionId === session.id || review.goalSessionId === session.id,
        ),
    ) ??
    workflow.items.find((item) => item.status !== "approved") ??
    workflow.items.at(-1);
  const review =
    active?.reviews.find((review) => review.summarySessionId === session.id) ??
    active?.reviews.at(-1);
  const reviewable =
    active?.status === "awaiting_review" &&
    review?.decision === "pending" &&
    review.id === active.reviews.at(-1)?.id;
  async function decide() {
    if (!active || !review || !reviewable) return;
    setBusy(true);
    setError("");
    try {
      const next = await decideWorkItem(active.id, review.id, feedback);
      setFeedback("");
      if (next) onSession(next.id);
      else await refresh();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not save review.");
    } finally {
      setBusy(false);
    }
  }
  async function resume(withReply = false) {
    if (!active) return;
    setBusy(true);
    setError("");
    try {
      const next =
        withReply && active.status === "blocked" && reply.trim()
          ? await replyToWorkItem(active.id, reply)
          : await retryImplementation(session.storyId);
      onSession(next.id);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not resume work.");
    } finally {
      setBusy(false);
    }
  }
  return (
    <section aria-label="Implementation" className="implementation-progress">
      <header className="section-toolbar">
        <h2 className="section-title">Work items</h2>
        <span>
          {workflow.items.filter((item) => item.status === "approved").length} /{" "}
          {workflow.items.length} approved
        </span>
      </header>
      <DataTable
        rows={workflow.items}
        columns={columns}
        rowKey={(item) => item.id}
        label="Work items"
        emptyMessage="No approved items"
      />
      {active && (
        <LayerCard className="work-item-review">
          <header className="section-toolbar">
            <h3 className="section-title">{active.title}</h3>
            <ItemBadge status={active.status} />
          </header>
          {active.blocker && <p className="error-message">{active.blocker}</p>}
          {active.status === "blocked" && (
            <>
              <InputArea
                label="Reply to blocker"
                value={reply}
                onChange={(event) => setReply(event.target.value)}
                rows={3}
                maxLength={8000}
                disabled={busy}
              />
              <div className="form-actions">
                <Button
                  variant="primary"
                  shape="square"
                  title="Resume work item"
                  aria-label="Resume work item"
                  disabled={(!reply.trim() && active.replies.length === 0) || busy}
                  loading={busy}
                  onClick={() => {
                    void resume(true);
                  }}
                >
                  <Play size={16} />
                </Button>
              </div>
            </>
          )}
          {(active.status === "failed" || active.status === "blocked") && (
            <div className="form-actions">
              <Button
                variant="secondary"
                shape="square"
                title="Retry work item"
                aria-label="Retry work item"
                disabled={busy}
                loading={busy}
                onClick={() => {
                  void resume();
                }}
              >
                <ArrowClockwise size={16} />
              </Button>
            </div>
          )}
          {review && (
            <section aria-label="Executive summary" className="executive-summary">
              <h4 className="section-title">Executive summary</h4>
              <div className="chat-markdown">
                <Markdown remarkPlugins={plugins}>{review.summary}</Markdown>
              </div>
              <h4 className="section-title">Verification</h4>
              <p>{review.verification}</p>
              {review.risks.length > 0 && (
                <>
                  <h4 className="section-title">Risks</h4>
                  <ul>
                    {review.risks.map((risk, index) => (
                      <li key={index}>{risk}</li>
                    ))}
                  </ul>
                </>
              )}
              {review.feedback && <blockquote>{review.feedback}</blockquote>}
              {reviewable && (
                <form
                  onSubmit={(event) => {
                    event.preventDefault();
                    void decide();
                  }}
                >
                  <InputArea
                    label="Feedback on work item"
                    rows={3}
                    value={feedback}
                    maxLength={8000}
                    disabled={busy}
                    onChange={(event) => setFeedback(event.target.value)}
                    onKeyDown={(event) => {
                      if (event.key === "Enter" && (event.metaKey || event.ctrlKey)) {
                        event.preventDefault();
                        event.currentTarget.form?.requestSubmit();
                      }
                    }}
                  />
                  <div className="form-actions">
                    <Button variant="primary" type="submit" disabled={busy} loading={busy}>
                      {feedback.trim() ? <ChatText size={16} /> : <Check size={16} />}{" "}
                      {feedback.trim() ? "Send Feedback" : "Approve"}
                    </Button>
                  </div>
                </form>
              )}
            </section>
          )}
          {(!review || review.goalSessionId === active.lastGoalSessionId) && (
            <details className="criterion-evidence">
              <summary>Acceptance criteria</summary>
              <ol>
                {active.acceptanceCriteria.map((criterion, index) => {
                  const evidence = active.criteria.find((entry) => entry.criterionIndex === index);
                  return (
                    <li key={index}>
                      <p>{criterion}</p>
                      {evidence && (
                        <>
                          <p>{evidence.summary}</p>
                          <ul>
                            {evidence.evidence.map((entry, index) => (
                              <li key={index}>{entry}</li>
                            ))}
                          </ul>
                        </>
                      )}
                    </li>
                  );
                })}
              </ol>
            </details>
          )}
          {error && (
            <p className="error-message" role="alert">
              {error}
            </p>
          )}
        </LayerCard>
      )}
    </section>
  );
}
