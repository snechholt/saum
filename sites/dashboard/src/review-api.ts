import { request } from "./api";
import {
  parsePlanComment,
  parseSession,
  type PlanComment,
  type PlanField,
  type Session,
} from "./sessions";
export type CommentDraft = {
  itemId: string;
  field: PlanField;
  criterionIndex: number;
  start: number;
  end: number;
  quote: string;
  body: string;
};
export async function addPlanComment(sessionId: string, input: CommentDraft): Promise<PlanComment> {
  return parsePlanComment(
    await request(`sessions/${encodeURIComponent(sessionId)}/comments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(input),
    }),
  );
}
export async function removePlanComment(sessionId: string, id: string): Promise<void> {
  await request(`sessions/${encodeURIComponent(sessionId)}/comments/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}
export async function approvePlan(sessionId: string): Promise<Session> {
  return parseSession(
    await request(`sessions/${encodeURIComponent(sessionId)}/approve`, { method: "POST" }),
  );
}
export async function sendPlanFeedback(sessionId: string): Promise<Session> {
  return parseSession(
    await request(`sessions/${encodeURIComponent(sessionId)}/feedback`, { method: "POST" }),
  );
}
