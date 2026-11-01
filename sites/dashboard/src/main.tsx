import { StrictMode } from "react";
import { BrowserRouter, Link, Navigate, Route, Routes } from "react-router";
import { createRoot } from "react-dom/client";
import "@cloudflare/kumo/styles/standalone";
import "@saum/ui/styles.css";
import "./styles.css";
import App from "./App";
import { LegacyStoryRedirect } from "./LegacyStoryRedirect";
import { applyTheme, initialTheme, SettingsProvider } from "./settings";
applyTheme(initialTheme());
const root = document.getElementById("root");
if (!root) throw new Error("Missing application root");
createRoot(root).render(
  <StrictMode>
    <BrowserRouter>
      <SettingsProvider>
        <Routes>
          <Route path="/" element={<Navigate to="/factories" replace />} />
          <Route path="/factories" element={<App />} />
          <Route path="/factories/:factoryId" element={<App />} />
          <Route path="/factories/:factoryId/stories/:storyId" element={<App />} />
          <Route
            path="/factories/:factoryId/stories/:storyId/sessions/:sessionId"
            element={<App />}
          />
          <Route path="/settings" element={<App />} />
          <Route path="/audit" element={<App />} />
          <Route path="/audit/:auditSessionId" element={<App />} />
          <Route path="/stories/:storyId" element={<LegacyStoryRedirect />} />
          <Route path="/stories/:storyId/sessions/:sessionId" element={<LegacyStoryRedirect />} />
          <Route
            path="*"
            element={
              <main className="workspace">
                <h1>Not found</h1>
                <Link to="/factories">Factories</Link>
              </main>
            }
          />
        </Routes>
      </SettingsProvider>
    </BrowserRouter>
  </StrictMode>,
);
