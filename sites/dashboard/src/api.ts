export type Factory = { id: string; name: string };
export type PlanningStatus =
  | "needs_goal"
  | "queued"
  | "preparing"
  | "running"
  | "completed"
  | "failed"
  | "interrupted"
  | "cancelled";
export type UserStory = {
  id: string;
  title: string;
  goal: string;
  planningStatus: PlanningStatus;
  sessionId: string;
};

export function parsePlanningStatus(value: unknown): PlanningStatus {
  switch (value) {
    case "needs_goal":
    case "queued":
    case "preparing":
    case "running":
    case "completed":
    case "failed":
    case "interrupted":
    case "cancelled":
      return value;
    default:
      throw new Error("Invalid planning status.");
  }
}
export type StoryDraft = Pick<UserStory, "title" | "goal">;

export function parseUserStory(value: unknown): UserStory {
  if (
    typeof value !== "object" ||
    value === null ||
    !("id" in value) ||
    typeof value.id !== "string" ||
    !value.id ||
    !("title" in value) ||
    typeof value.title !== "string" ||
    !value.title ||
    !("goal" in value) ||
    typeof value.goal !== "string"
  ) {
    throw new Error("The server returned an invalid user story.");
  }
  if (
    !("planningStatus" in value) ||
    !("sessionId" in value) ||
    typeof value.sessionId !== "string"
  )
    throw new Error("Invalid story session.");
  return {
    id: value.id,
    title: value.title,
    goal: value.goal,
    planningStatus: parsePlanningStatus(value.planningStatus),
    sessionId: value.sessionId,
  };
}

type NamedItem = { id: string; name: string };
export function parseNamedItem(value: unknown): NamedItem {
  if (
    typeof value !== "object" ||
    value === null ||
    !("id" in value) ||
    typeof value.id !== "string" ||
    !value.id ||
    !("name" in value) ||
    typeof value.name !== "string" ||
    !value.name
  ) {
    throw new Error("The server returned an invalid response.");
  }
  return { id: value.id, name: value.name };
}
export async function request(path: string, init?: RequestInit): Promise<unknown> {
  const response = await fetch(`/api/${path}`, init);
  if (response.status === 204) return null;
  const body: unknown = await response.json();
  if (!response.ok) {
    const message =
      typeof body === "object" && body !== null && "error" in body && typeof body.error === "string"
        ? body.error
        : "Request failed.";
    throw new Error(message);
  }
  return body;
}
function parseList(body: unknown): NamedItem[] {
  if (!Array.isArray(body)) throw new Error("The server returned an invalid list.");
  return body.map((value: unknown) => parseNamedItem(value));
}
export async function listFactories(): Promise<Factory[]> {
  return parseList(await request("factories"));
}
export async function createFactory(name: string): Promise<Factory> {
  return parseNamedItem(
    await request("factories", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    }),
  );
}
export async function updateFactory(id: Factory["id"], name: string): Promise<Factory> {
  return parseNamedItem(
    await request(`factories/${encodeURIComponent(id)}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    }),
  );
}
export async function deleteFactory(id: Factory["id"]): Promise<void> {
  await request(`factories/${encodeURIComponent(id)}`, { method: "DELETE" });
}
export async function listStories(factoryID: Factory["id"]): Promise<UserStory[]> {
  const body = await request(`factories/${encodeURIComponent(factoryID)}/stories`);
  if (!Array.isArray(body)) throw new Error("The server returned an invalid story list.");
  return body.map((value: unknown) => parseUserStory(value));
}
export async function createStory(factoryID: Factory["id"], draft: StoryDraft): Promise<UserStory> {
  return parseUserStory(
    await request(`factories/${encodeURIComponent(factoryID)}/stories`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(draft),
    }),
  );
}
