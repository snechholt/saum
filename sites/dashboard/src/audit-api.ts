import {
  parseFactory,
  parsePlanningStatus,
  parseSessionKind,
  request,
  type Factory,
  type SessionKind,
} from "./api";
import type { Session } from "./sessions";
export type AuditSession = Pick<
  Session,
  "id" | "storyId" | "status" | "kind" | "createdAt" | "finishedAt" | "threadId" | "error"
>;
function parseAuditSession(value: unknown): AuditSession {
  if (typeof value !== "object" || value === null || Array.isArray(value))
    throw new Error("Invalid audit session.");
  const row = value as Record<string, unknown>;
  function text(key: string): string {
    const value = row[key];
    if (typeof value !== "string") throw new Error(`Invalid session field: ${key}`);
    return value;
  }
  return {
    id: text("id"),
    storyId: text("storyId"),
    status: parsePlanningStatus(row["status"]),
    kind: parseSessionKind(row["kind"]),
    createdAt: text("createdAt"),
    finishedAt: text("finishedAt"),
    threadId: text("threadId"),
    error: text("error"),
  };
}
export type AuditEntry = {
  cursor: number;
  session: AuditSession;
  factory: Factory;
  storyTitle: string;
  workItemTitle: string;
};
export type AuditPage = { entries: AuditEntry[]; nextCursor: number };
export async function getAudit(
  factoryId: string,
  role: SessionKind | "all",
  before = 0,
): Promise<AuditPage> {
  const query = new URLSearchParams();
  if (factoryId !== "all") query.set("factory", factoryId);
  if (role !== "all") query.set("role", role);
  if (before) query.set("before", String(before));
  const result = await request(`audit?${query}`);
  if (
    typeof result !== "object" ||
    result === null ||
    !("entries" in result) ||
    !Array.isArray(result.entries) ||
    !("nextCursor" in result) ||
    typeof result.nextCursor !== "number" ||
    !Number.isSafeInteger(result.nextCursor) ||
    result.nextCursor < 0
  )
    throw new Error("Invalid audit response.");
  return {
    nextCursor: result.nextCursor,
    entries: result.entries.map((value: unknown) => {
      if (
        typeof value !== "object" ||
        value === null ||
        !("cursor" in value) ||
        typeof value.cursor !== "number" ||
        !Number.isSafeInteger(value.cursor) ||
        !("session" in value) ||
        !("factory" in value) ||
        !("storyTitle" in value) ||
        typeof value.storyTitle !== "string" ||
        !("workItemTitle" in value) ||
        typeof value.workItemTitle !== "string"
      )
        throw new Error("Invalid audit entry.");
      return {
        cursor: value.cursor,
        session: parseAuditSession(value.session),
        factory: parseFactory(value.factory),
        storyTitle: value.storyTitle,
        workItemTitle: value.workItemTitle,
      };
    }),
  };
}
