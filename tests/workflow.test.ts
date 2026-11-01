import { expect, test } from "bun:test";
import { parseWorkflow, parseSession } from "../sites/dashboard/src/sessions";
import { parseSettings } from "../sites/dashboard/src/settings-api";
import { sessionPath, storyPath } from "../sites/dashboard/src/paths";
const item = {
  id: "item",
  sessionId: "plan",
  key: "login",
  title: "Sign in",
  description: "Implement sign-in",
  acceptanceCriteria: ["Valid login works"],
  planSessionId: "plan",
  position: 0,
  status: "awaiting_review",
  blocker: "",
  currentSessionId: "summary",
  lastGoalSessionId: "goal",
  criteria: [{ criterionIndex: 0, summary: "Verified", evidence: ["auth tests pass"] }],
  reviews: [
    {
      id: "review",
      workItemId: "item",
      goalSessionId: "goal",
      summarySessionId: "summary",
      summary: "Users can sign in.",
      verification: "Observed passing auth tests.",
      risks: [],
      decision: "pending",
      feedback: "",
      createdAt: "now",
      reviewedAt: "",
    },
  ],
  replies: [],
};
test("parses pending executive review without treating it as approval", () => {
  const workflow = parseWorkflow({ planSessionId: "plan", items: [item] });
  expect(workflow?.items[0]?.status).toBe("awaiting_review");
  expect(workflow?.items[0]?.reviews[0]?.decision).toBe("pending");
});
test("rejects unknown work item states and invalid criterion indices", () => {
  expect(() =>
    parseWorkflow({ planSessionId: "plan", items: [{ ...item, status: "shipped" }] }),
  ).toThrow();
  expect(() =>
    parseWorkflow({
      planSessionId: "plan",
      items: [{ ...item, criteria: [{ criterionIndex: -1, summary: "Wrong", evidence: [] }] }],
    }),
  ).toThrow();
});
test("factory agents have an explicit summary session kind", () => {
  const session = parseSession({
    id: "summary",
    storyId: "story",
    status: "completed",
    kind: "summary",
    worktree: "/factory/main",
    threadId: "factory-thread",
    prompt: "review",
    error: "",
    planSubmitted: false,
    createdAt: "now",
    finishedAt: "now",
    parentSessionId: "plan",
    feedbackBatchId: "",
    reviewStatus: "pending",
    reviewedAt: "",
  });
  expect(session.kind).toBe("summary");
  expect(session.planSubmitted).toBe(false);
});
test("settings accept only supported themes and a directory string", () => {
  expect(parseSettings({ factoriesRoot: "/factories", theme: "system" }).theme).toBe("system");
  expect(() => parseSettings({ factoriesRoot: "/factories", theme: "automatic" })).toThrow();
  expect(() => parseSettings({ theme: "dark" })).toThrow();
});
test("all story and session links are factory-scoped and escape segments", () => {
  expect(storyPath("factory", "story")).toBe("/factories/factory/stories/story");
  expect(sessionPath("factory", "story", "session")).toBe(
    "/factories/factory/stories/story/sessions/session",
  );
  expect(storyPath("a/b", "x y")).toBe("/factories/a%2Fb/stories/x%20y");
});
