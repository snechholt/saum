export type Project = { id: string; name: string };

function parseProject(value: unknown): Project {
  if (
    typeof value !== "object" ||
    value === null ||
    !("id" in value) ||
    typeof value.id !== "string" ||
    !("name" in value) ||
    typeof value.name !== "string"
  ) {
    throw new Error("The server returned an invalid project.");
  }
  return { id: value.id, name: value.name };
}
async function request(init?: RequestInit): Promise<unknown> {
  const response = await fetch("/api/projects", init);
  const body: unknown = await response.json();
  if (!response.ok) {
    const message =
      typeof body === "object" && body !== null && "error" in body && typeof body.error === "string"
        ? body.error
        : "Could not reach the factory.";
    throw new Error(message);
  }
  return body;
}
export async function listProjects(): Promise<Project[]> {
  const body = await request();
  if (!Array.isArray(body)) throw new Error("The server returned an invalid project list.");
  return body.map((value: unknown) => parseProject(value));
}
export async function createProject(name: string): Promise<Project> {
  return parseProject(
    await request({
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    }),
  );
}
