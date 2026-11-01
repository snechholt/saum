import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { getSettings, putSettings, type Settings, type Theme } from "./settings-api";
export function initialTheme(): Theme {
  try {
    const value = localStorage.getItem("saum.theme");
    if (value === "dark" || value === "light") return value;
  } catch {}
  return "system";
}
export function applyTheme(theme: Theme) {
  const mode =
    theme === "system"
      ? matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      : theme;
  document.documentElement.dataset["mode"] = mode;
  document.documentElement.style.colorScheme = mode;
}
type SettingsContextValue = {
  settings: Settings | null;
  error: string;
  save: (settings: Settings) => Promise<void>;
};
const SettingsContext = createContext<SettingsContextValue | null>(null);
export function SettingsProvider({ children }: { children: ReactNode }) {
  const [settings, setSettings] = useState<Settings | null>(null),
    [error, setError] = useState("");
  useEffect(() => {
    let ignore = false;
    getSettings()
      .then((value) => {
        if (!ignore) setSettings(value);
      })
      .catch((reason) => {
        if (!ignore)
          setError(reason instanceof Error ? reason.message : "Could not load settings.");
      });
    return () => {
      ignore = true;
    };
  }, []);
  const theme = settings?.theme ?? initialTheme();
  useEffect(() => {
    const system = matchMedia("(prefers-color-scheme: dark)"),
      update = () => applyTheme(theme);
    update();
    system.addEventListener("change", update);
    try {
      localStorage.setItem("saum.theme", theme);
    } catch {}
    return () => system.removeEventListener("change", update);
  }, [theme]);
  async function save(draft: Settings) {
    const settings = await putSettings(draft);
    setSettings(settings);
    setError("");
  }
  return (
    <SettingsContext.Provider value={{ settings, error, save }}>
      {children}
    </SettingsContext.Provider>
  );
}
export function useSettings() {
  const value = useContext(SettingsContext);
  if (!value) throw new Error("Missing settings provider");
  return value;
}
