import { request } from "./api";
import { parseSession, type Session } from "./sessions";
export async function decideWorkItem(
  itemId: string,
  reviewId: string,
  feedback: string,
): Promise<Session | null> {
  const value = await request(
    `work-items/${encodeURIComponent(itemId)}/${feedback.trim() ? "feedback" : "approve"}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ reviewId, feedback }),
    },
  );
  if (typeof value !== "object" || value === null || !("session" in value))
    throw new Error("Invalid review response.");
  return value.session === null ? null : parseSession(value.session);
}
export async function replyToWorkItem(itemId: string, reply: string): Promise<Session> {
  return parseSession(
    await request(`work-items/${encodeURIComponent(itemId)}/reply`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ reply }),
    }),
  );
}
