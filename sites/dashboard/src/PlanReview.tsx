import { useCallback, useMemo, useRef, useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { InputArea } from "@cloudflare/kumo/components/input";
import { Popover } from "@cloudflare/kumo/components/popover";
import { ChatText, Check, Plus, Trash, X } from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import {
  addPlanComment,
  approvePlan,
  removePlanComment,
  sendPlanFeedback,
  type CommentDraft,
} from "./review-api";
import type { ActionItem, PlanComment, PlanField, Session } from "./sessions";
import { annotatedSegments } from "./plan-annotations";

type SelectedText = Omit<CommentDraft, "body"> & { range: Range; element: HTMLElement };
function PlanText({
  item,
  field,
  criterionIndex = -1,
  text,
  comments,
}: {
  item: ActionItem;
  field: PlanField;
  criterionIndex?: number;
  text: string;
  comments: PlanComment[];
}) {
  const anchors = comments.filter(
    (comment) =>
      comment.itemId === item.id &&
      comment.field === field &&
      comment.criterionIndex === criterionIndex,
  );
  return (
    <span
      className="plan-text"
      data-plan-item={item.id}
      data-plan-field={field}
      data-criterion-index={criterionIndex}
    >
      {annotatedSegments(text, anchors).map((segment) =>
        segment.comments.length ? (
          <mark
            key={segment.start}
            title={segment.comments.map((comment) => comment.body).join("\n\n")}
          >
            {segment.text}
          </mark>
        ) : (
          <span key={segment.start}>{segment.text}</span>
        ),
      )}
    </span>
  );
}
export function PlanReview({
  session,
  items,
  comments,
  latest,
  refresh,
  onRevision,
}: {
  session: Session | null;
  items: ActionItem[];
  comments: PlanComment[];
  latest: boolean;
  refresh: () => Promise<void>;
  onRevision: (session: Session) => void;
}) {
  const root = useRef<HTMLDivElement>(null);
  const [selection, setSelection] = useState<SelectedText | null>(null),
    [body, setBody] = useState(""),
    [busy, setBusy] = useState(false),
    [error, setError] = useState("");
  const focusComment = useCallback((element: HTMLTextAreaElement | null) => element?.focus(), []);
  const editable =
    !!session &&
    session.status === "completed" &&
    session.planSubmitted &&
    session.reviewStatus === "pending" &&
    latest;
  const drafts = comments.filter((comment) => !comment.batchId);
  const anchor = useMemo(
    () => ({
      getBoundingClientRect: () => selection?.range.getBoundingClientRect() ?? new DOMRect(),
      ...(selection ? { contextElement: selection.element } : {}),
    }),
    [selection],
  );
  function selectText() {
    if (!editable || busy) return;
    const selected = window.getSelection();
    if (!selected || selected.rangeCount === 0 || selected.isCollapsed) return;
    const range = selected.getRangeAt(0);
    const quote = range.toString();
    if (!quote.trim()) return;
    const startNode = range.startContainer;
    const startElement = startNode instanceof Element ? startNode : startNode.parentElement;
    const element = startElement?.closest<HTMLElement>("[data-plan-field]");
    if (!element || !root.current?.contains(element) || !element.contains(range.endContainer))
      return;
    const itemId = element.dataset["planItem"],
      field = element.dataset["planField"],
      criterionIndex = Number(element.dataset["criterionIndex"]);
    if (
      !itemId ||
      !items.some((item) => item.id === itemId) ||
      (field !== "title" && field !== "description" && field !== "acceptanceCriteria") ||
      !Number.isSafeInteger(criterionIndex)
    )
      return;
    const before = range.cloneRange();
    before.selectNodeContents(element);
    before.setEnd(range.startContainer, range.startOffset);
    const start = before.toString().length,
      end = start + quote.length;
    if (element.textContent?.slice(start, end) !== quote) return;
    setBody("");
    setError("");
    setSelection({
      itemId,
      field,
      criterionIndex,
      start,
      end,
      quote,
      range: range.cloneRange(),
      element,
    });
  }
  async function saveComment() {
    if (!session || !selection || !body.trim() || busy) return;
    setBusy(true);
    setError("");
    try {
      const { range: _range, element: _element, ...selected } = selection;
      await addPlanComment(session.id, { ...selected, body });
      await refresh();
      setSelection(null);
      setBody("");
      window.getSelection()?.removeAllRanges();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not save comment.");
    } finally {
      setBusy(false);
    }
  }
  async function removeComment(comment: PlanComment) {
    if (!session || busy) return;
    setBusy(true);
    setError("");
    try {
      await removePlanComment(session.id, comment.id);
      await refresh();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not remove comment.");
    } finally {
      setBusy(false);
    }
  }
  async function submitReview() {
    if (!session || !editable || busy) return;
    setBusy(true);
    setError("");
    try {
      if (drafts.length) {
        const next = await sendPlanFeedback(session.id);
        onRevision(next);
      } else {
        const implementation = await approvePlan(session.id);
        onRevision(implementation);
      }
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not review plan.");
    } finally {
      setBusy(false);
    }
  }
  const columns: Column<ActionItem>[] = [
    {
      id: "title",
      header: "Item",
      value: (item) => item.title,
      render: (item) => (
        <PlanText item={item} field="title" text={item.title} comments={comments} />
      ),
    },
    {
      id: "description",
      header: "Work",
      value: (item) => item.description,
      render: (item) => (
        <PlanText item={item} field="description" text={item.description} comments={comments} />
      ),
    },
    {
      id: "criteria",
      header: "Acceptance criteria",
      value: (item) => item.acceptanceCriteria.join(" · "),
      render: (item) => (
        <ul className="item-criteria">
          {item.acceptanceCriteria.map((criterion, index) => (
            <li key={index}>
              <PlanText
                item={item}
                field="acceptanceCriteria"
                criterionIndex={index}
                text={criterion}
                comments={comments}
              />
            </li>
          ))}
        </ul>
      ),
    },
  ];
  const label =
    session?.reviewStatus === "approved"
      ? "Approved"
      : session?.reviewStatus === "changes_requested"
        ? "Feedback sent"
        : drafts.length
          ? "Send Feedback"
          : "Approve";
  return (
    <div className="plan-review" ref={root}>
      <header className="section-toolbar">
        <h2 className="section-title">Plan</h2>
        <Button
          variant="primary"
          disabled={!editable || busy || selection !== null}
          loading={busy}
          onClick={() => {
            void submitReview();
          }}
        >
          {drafts.length ? <ChatText size={16} /> : <Check size={16} />} {label}
        </Button>
      </header>
      {error && !selection && (
        <p className="error-message" role="alert">
          {error}
        </p>
      )}
      {/* oxlint-disable-next-line jsx-a11y/no-noninteractive-element-interactions -- Observe native text selection without misrepresenting the readable table as a button. */}
      <section
        aria-label="Plan text"
        onMouseUp={selectText}
        onKeyUp={selectText}
        className={editable ? "plan-annotatable" : ""}
      >
        <DataTable
          rows={items}
          columns={columns}
          rowKey={(item) => item.id}
          label="Plan"
          emptyMessage={
            session?.status === "running" ||
            session?.status === "preparing" ||
            session?.status === "queued"
              ? "Planning…"
              : "No plan items submitted"
          }
        />
      </section>
      {comments.length > 0 && (
        <section className="plan-comments" aria-label="Plan comments">
          {comments.map((comment) => (
            <article className="plan-comment" key={comment.id}>
              <div className="plan-comment-heading">
                <span>
                  {comment.field === "title"
                    ? "Item"
                    : comment.field === "description"
                      ? "Work"
                      : `Acceptance criterion ${comment.criterionIndex + 1}`}
                </span>
                {editable && !comment.batchId && (
                  <Button
                    variant="ghost"
                    size="sm"
                    shape="square"
                    aria-label={`Remove comment ${comment.id}`}
                    title="Remove comment"
                    disabled={busy}
                    onClick={() => {
                      void removeComment(comment);
                    }}
                  >
                    <Trash size={16} />
                  </Button>
                )}
              </div>
              <blockquote>{comment.quote}</blockquote>
              <p>{comment.body}</p>
            </article>
          ))}
        </section>
      )}
      <Popover.Root
        open={selection !== null}
        onOpenChange={(open) => {
          if (!open && !busy) setSelection(null);
        }}
      >
        <Popover.Portal>
          <Popover.Positioner
            anchor={anchor}
            side="bottom"
            align="start"
            sideOffset={8}
            positionMethod="fixed"
          >
            <Popover.Popup className="feedback-popup">
              <Popover.Title>Add comment</Popover.Title>
              <blockquote>{selection?.quote}</blockquote>
              <form
                onSubmit={(event) => {
                  event.preventDefault();
                  void saveComment();
                }}
              >
                <InputArea
                  label="Comment"
                  value={body}
                  onChange={(event) => setBody(event.target.value)}
                  rows={3}
                  maxLength={4000}
                  disabled={busy}
                  ref={focusComment}
                  onKeyDown={(event) => {
                    if (event.key === "Enter" && (event.metaKey || event.ctrlKey)) {
                      event.preventDefault();
                      event.currentTarget.form?.requestSubmit();
                    }
                  }}
                />
                {error && (
                  <p className="error-message" role="alert">
                    {error}
                  </p>
                )}
                <div className="form-actions">
                  <Button
                    variant="ghost"
                    shape="square"
                    type="button"
                    aria-label="Cancel comment"
                    title="Cancel"
                    disabled={busy}
                    onClick={() => setSelection(null)}
                  >
                    <X size={16} />
                  </Button>
                  <Button
                    variant="primary"
                    shape="square"
                    type="submit"
                    aria-label="Add comment"
                    title="Add comment"
                    disabled={!body.trim() || busy}
                    loading={busy}
                  >
                    <Plus size={16} />
                  </Button>
                </div>
              </form>
            </Popover.Popup>
          </Popover.Positioner>
        </Popover.Portal>
      </Popover.Root>
    </div>
  );
}
