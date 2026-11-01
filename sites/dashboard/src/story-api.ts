import { parseFactory, parseUserStory, request, type Factory, type UserStory } from "./api";
import { parseSession, type Session } from "./sessions";
export type StoryDetail = { story: UserStory; factory: Factory; sessions: Session[] };
export async function getStory(id: string): Promise<StoryDetail> {
  const body = await request(`stories/${encodeURIComponent(id)}`);
  if (
    typeof body !== "object" ||
    body === null ||
    !("story" in body) ||
    !("factory" in body) ||
    !("sessions" in body) ||
    !Array.isArray(body.sessions)
  )
    throw new Error("Invalid story response.");
  return {
    story: parseUserStory(body.story),
    factory: parseFactory(body.factory),
    sessions: body.sessions.map((value: unknown) => parseSession(value)),
  };
}
