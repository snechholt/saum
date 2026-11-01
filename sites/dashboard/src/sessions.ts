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
export type ImplementationItem = {
  key: string;
  title: string;
  status: "pending" | "running" | "completed" | "blocked";
  summary: string;
};
export type Implementation = {
  planSessionId: string;
  submitted: boolean;
  summary: string;
  items: ImplementationItem[];
};
export type SessionDetail = {
  planSession: Session;
  implementation: Implementation | null;
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
export function parseImplementation(value: unknown): Implementation | null {
  if (value === null) return null;
  const row = record(value),
    items = row["items"];
  if (typeof row["submitted"] !== "boolean" || !Array.isArray(items))
    throw new Error("Invalid implementation.");
  return {
    planSessionId: text(row, "planSessionId"),
    submitted: row["submitted"],
    summary: text(row, "summary"),
    items: items.map((value: unknown) => {
      const item = record(value),
        status = item["status"];
      if (
        status !== "pending" &&
        status !== "running" &&
        status !== "completed" &&
        status !== "blocked"
      )
        throw new Error("Invalid item progress.");
      return {
        key: text(item, "key"),
        title: text(item, "title"),
        summary: text(item, "summary"),
        status,
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
    implementation = parseImplementation(row["implementation"]);
  if (
    planSession.kind !== "planning" ||
    (session.kind === "implementation"
      ? !implementation ||
        planSession.id !== implementation.planSessionId ||
        planSession.reviewStatus !== "approved"
      : implementation !== null || planSession.id !== session.id)
  )
    throw new Error("Invalid plan context.");
  return {
    session,
    planSession,
    implementation,
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
    if (status === "completed") return "Built";
    if (status === "running") return "Building";
    if (status === "queued") return "Build queued";
  }
  return status === "completed"
    ? "Planned"
    : status === "needs_goal"
      ? "Needs goal"
      : status[0]?.toUpperCase() + status.slice(1);
}
