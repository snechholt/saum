import { request } from "./api";
export type Theme = "system" | "light" | "dark";
export type Settings = { factoriesRoot: string; theme: Theme };
export function parseSettings(value: unknown): Settings {
  if (
    typeof value !== "object" ||
    value === null ||
    !("factoriesRoot" in value) ||
    typeof value.factoriesRoot !== "string" ||
    !("theme" in value) ||
    (value.theme !== "system" && value.theme !== "light" && value.theme !== "dark")
  )
    throw new Error("Invalid settings response.");
  return { factoriesRoot: value.factoriesRoot, theme: value.theme };
}

export async function getSettings(): Promise<Settings> {
  return parseSettings(await request("settings"));
}
export async function putSettings(settings: Settings): Promise<Settings> {
  return parseSettings(
    await request("settings", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(settings),
    }),
  );
}
