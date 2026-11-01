export type Factory = { id: string; name: string };
export type Project = { id: string; name: string };

type NamedItem = { id: string; name: string };
function parseNamedItem(value: unknown): NamedItem {
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
async function request(path: string, init?: RequestInit): Promise<unknown> {
  const response = await fetch(`/api/${path}`, init);
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
export async function listProjects(factoryID: Factory["id"]): Promise<Project[]> {
  return parseList(await request(`factories/${encodeURIComponent(factoryID)}/projects`));
}
