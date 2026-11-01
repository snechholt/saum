import { StrictMode } from "react";
import { BrowserRouter, Link, Route, Routes } from "react-router";
import { createRoot } from "react-dom/client";
import "@cloudflare/kumo/styles/standalone";
import "@saum/ui/styles.css";
import "./styles.css";
import App from "./App";

const root = document.getElementById("root");
if (!root) throw new Error("Missing application root");
createRoot(root).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/stories/:storyId" element={<App />} />
        <Route path="/stories/:storyId/sessions/:sessionId" element={<App />} />
        <Route path="/factories/:factoryId" element={<App />} />
        <Route
          path="*"
          element={
            <main className="workspace">
              <h1>Not found</h1>
              <Link to="/">Factories</Link>
            </main>
          }
        />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);
