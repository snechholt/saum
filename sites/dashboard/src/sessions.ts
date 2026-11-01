import {
  parsePlanningStatus,
  parseSessionKind,
  request,
  type PlanningStatus,
  type SessionKind,
} from "./api";

export type ReviewStatus = "pending" | "approved" | "changes_requested";
export type PlanField = "title" | "description" | "acceptanceCriteria";
export type PlanComment = {
  id: string;
  sessionId: string;
  itemId: string;
  field: PlanField;
  criterionIndex: number;
  start: number;
  end: number;
  quote: string;
  body: string;
  batchId: string;
  createdAt: string;
};

export type Session = {
  id: string;
  storyId: string;
  status: PlanningStatus;
  kind: SessionKind;
  worktree: string;
  threadId: string;
  prompt: string;
  error: string;
  planSubmitted: boolean;
  createdAt: string;
  finishedAt: string;
  parentSessionId: string;
  feedbackBatchId: string;
  reviewStatus: ReviewStatus;
  reviewedAt: string;
};
export type ActionItem = {
  id: string;
  sessionId: string;
  key: string;
  title: string;
  description: string;
  acceptanceCriteria: string[];
};
export type WorkItemStatus =
  | "pending"
  | "queued"
  | "working"
  | "blocked"
  | "review_requested"
  | "summarizing"
  | "awaiting_review"
  | "approved"
  | "failed";
export type CriterionSubmission = { criterionIndex: number; summary: string; evidence: string[] };
export type ExecutiveReview = {
  id: string;
  workItemId: string;
  goalSessionId: string;
  summarySessionId: string;
  summary: string;
  verification: string;
  risks: string[];
  decision: ReviewStatus;
  feedback: string;
  createdAt: string;
  reviewedAt: string;
};
export type WorkItem = ActionItem & {
  planSessionId: string;
  position: number;
  status: WorkItemStatus;
  blocker: string;
  currentSessionId: string;
  lastGoalSessionId: string;
  criteria: CriterionSubmission[];
  reviews: ExecutiveReview[];
  replies: string[];
};
export type GoalWorkflow = { planSessionId: string; items: WorkItem[] };
export type SessionDetail = {
  planSession: Session;
  workflow: GoalWorkflow | null;
  session: Session;
  items: ActionItem[];
  comments: PlanComment[];
  receivedFeedback: PlanComment[];
};
export type SessionEvent = { id: number; stream: string; data: string };

function record(value: unknown): Record<string, unknown> {
  if (typeof value !== "object" || value === null || Array.isArray(value))
    throw new Error("Invalid session response.");
  return value as Record<string, unknown>;
}
function text(value: Record<string, unknown>, key: string): string {
  const result = value[key];
  if (typeof result !== "string") throw new Error(`Invalid session field: ${key}`);
  return result;
}
function reviewStatus(value: unknown): ReviewStatus {
  if (value === "pending" || value === "approved" || value === "changes_requested") return value;
  throw new Error("Invalid plan review status.");
}
export function parsePlanComment(value: unknown): PlanComment {
  const row = record(value),
    field = row["field"];
  if (field !== "title" && field !== "description" && field !== "acceptanceCriteria")
    throw new Error("Invalid feedback field.");
  const index = row["criterionIndex"],
    start = row["start"],
    end = row["end"];
  if (
    typeof index !== "number" ||
    !Number.isSafeInteger(index) ||
    typeof start !== "number" ||
    !Number.isSafeInteger(start) ||
    typeof end !== "number" ||
    !Number.isSafeInteger(end)
  )
    throw new Error("Invalid comment anchor.");
  return {
    id: text(row, "id"),
    sessionId: text(row, "sessionId"),
    itemId: text(row, "itemId"),
    field,
    criterionIndex: index,
    start,
    end,
    quote: text(row, "quote"),
    body: text(row, "body"),
    batchId: text(row, "batchId"),
    createdAt: text(row, "createdAt"),
  };
}
export function parseSession(value: unknown): Session {
  const row = record(value);
  if (typeof row["planSubmitted"] !== "boolean") throw new Error("Invalid plan submission state.");
  return {
    id: text(row, "id"),
    storyId: text(row, "storyId"),
    status: parsePlanningStatus(row["status"]),
    kind: parseSessionKind(row["kind"]),
    worktree: text(row, "worktree"),
    threadId: text(row, "threadId"),
    prompt: text(row, "prompt"),
    error: text(row, "error"),
    planSubmitted: row["planSubmitted"],
    createdAt: text(row, "createdAt"),
    finishedAt: text(row, "finishedAt"),
    parentSessionId: text(row, "parentSessionId"),
    feedbackBatchId: text(row, "feedbackBatchId"),
    reviewStatus: reviewStatus(row["reviewStatus"]),
    reviewedAt: text(row, "reviewedAt"),
  };
}
function parseItem(value: unknown): ActionItem {
  const row = record(value),
    criteria = row["acceptanceCriteria"];
  if (!Array.isArray(criteria)) throw new Error("Invalid acceptance criteria.");
  return {
    id: text(row, "id"),
    sessionId: text(row, "sessionId"),
    key: text(row, "key"),
    title: text(row, "title"),
    description: text(row, "description"),
    acceptanceCriteria: criteria.map((criterion: unknown) => {
      if (typeof criterion !== "string") throw new Error("Invalid criterion.");
      return criterion;
    }),
  };
}
function strings(value: unknown): string[] {
  if (!Array.isArray(value)) throw new Error("Expected a list.");
  return value.map((entry: unknown) => {
    if (typeof entry !== "string") throw new Error("Expected text.");
    return entry;
  });
}
function integer(value: unknown): number {
  if (typeof value !== "number" || !Number.isSafeInteger(value) || value < 0)
    throw new Error("Invalid index.");
  return value;
}
function workItemStatus(value: unknown): WorkItemStatus {
  switch (value) {
    case "pending":
    case "queued":
    case "working":
    case "blocked":
    case "review_requested":
    case "summarizing":
    case "awaiting_review":
    case "approved":
    case "failed":
      return value;
    default:
      throw new Error("Invalid work item status.");
  }
}
function parseExecutiveReview(value: unknown): ExecutiveReview {
  const row = record(value);
  return {
    id: text(row, "id"),
    workItemId: text(row, "workItemId"),
    goalSessionId: text(row, "goalSessionId"),
    summarySessionId: text(row, "summarySessionId"),
    summary: text(row, "summary"),
    verification: text(row, "verification"),
    risks: strings(row["risks"]),
    decision: reviewStatus(row["decision"]),
    feedback: text(row, "feedback"),
    createdAt: text(row, "createdAt"),
    reviewedAt: text(row, "reviewedAt"),
  };
}
export function parseWorkflow(value: unknown): GoalWorkflow | null {
  if (value === null) return null;
  const row = record(value),
    items = row["items"];
  if (!Array.isArray(items)) throw new Error("Invalid goal workflow.");
  return {
    planSessionId: text(row, "planSessionId"),
    items: items.map((value: unknown) => {
      const item = record(value),
        criteria = item["criteria"],
        reviews = item["reviews"];
      if (!Array.isArray(criteria) || !Array.isArray(reviews))
        throw new Error("Invalid work item review.");
      return {
        ...parseItem(item),
        planSessionId: text(item, "planSessionId"),
        position: integer(item["position"]),
        status: workItemStatus(item["status"]),
        blocker: text(item, "blocker"),
        currentSessionId: text(item, "currentSessionId"),
        lastGoalSessionId: text(item, "lastGoalSessionId"),
        replies: strings(item["replies"]),
        reviews: reviews.map((value: unknown) => parseExecutiveReview(value)),
        criteria: criteria.map((value: unknown) => {
          const criterion = record(value);
          return {
            criterionIndex: integer(criterion["criterionIndex"]),
            summary: text(criterion, "summary"),
            evidence: strings(criterion["evidence"]),
          };
        }),
      };
    }),
  };
}
export async function getSession(id: string): Promise<SessionDetail> {
  const row = record(await request(`sessions/${encodeURIComponent(id)}`)),
    items = row["items"];
  if (!Array.isArray(items)) throw new Error("Invalid session items.");
  const comments = row["comments"],
    received = row["receivedFeedback"];
  if (!Array.isArray(comments) || !Array.isArray(received))
    throw new Error("Invalid plan feedback.");
  const session = parseSession(row["session"]),
    planSession = parseSession(row["planSession"]),
    workflow = parseWorkflow(row["workflow"]);
  if (
    planSession.kind !== "planning" ||
    (session.kind !== "planning"
      ? !workflow ||
        planSession.id !== workflow.planSessionId ||
        planSession.reviewStatus !== "approved"
      : planSession.id !== session.id)
  )
    throw new Error("Invalid plan context.");
  return {
    session,
    planSession,
    workflow,
    items: items.map((item: unknown) => parseItem(item)),
    comments: comments.map((comment: unknown) => parsePlanComment(comment)),
    receivedFeedback: received.map((comment: unknown) => parsePlanComment(comment)),
  };
}
export async function getSessionEvents(id: string, after: number): Promise<SessionEvent[]> {
  const rows = await request(`sessions/${encodeURIComponent(id)}/events?after=${after}`);
  if (!Array.isArray(rows)) throw new Error("Invalid session events.");
  return rows.map((value: unknown) => {
    const row = record(value),
      id = row["id"];
    if (typeof id !== "number" || !Number.isSafeInteger(id) || id < 1)
      throw new Error("Invalid event ID.");
    return { id, stream: text(row, "stream"), data: text(row, "data") };
  });
}
export async function retryPlanning(storyId: string): Promise<Session> {
  return parseSession(
    await request(`stories/${encodeURIComponent(storyId)}/planning`, { method: "POST" }),
  );
}
export async function retryImplementation(storyId: string): Promise<Session> {
  return parseSession(
    await request(`stories/${encodeURIComponent(storyId)}/implementation`, { method: "POST" }),
  );
}
export async function cancelSession(id: string): Promise<void> {
  await request(`sessions/${encodeURIComponent(id)}`, { method: "DELETE" });
}
export function isActive(status: PlanningStatus): boolean {
  return status === "queued" || status === "preparing" || status === "running";
}
export function planningLabel(status: PlanningStatus, kind: SessionKind = "planning"): string {
  if (kind === "implementation") {
    if (status === "completed") return "Review requested";
    if (status === "running") return "Building";
    if (status === "queued") return "Build queued";
  }
  if (kind === "summary") {
    if (status === "completed") return "Summary ready";
    if (status === "running") return "Summarizing";
    if (status === "queued") return "Summary queued";
  }
  return status === "completed"
    ? "Planned"
    : status === "needs_goal"
      ? "Needs goal"
      : status[0]?.toUpperCase() + status.slice(1);
}
