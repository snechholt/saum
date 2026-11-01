import { expect, test } from "bun:test";
import { annotatedSegments } from "../sites/dashboard/src/plan-annotations";
import type { PlanComment } from "../sites/dashboard/src/sessions";
function comment(start: number, end: number, quote: string, id = "comment"): PlanComment {
  return {
    id,
    sessionId: "plan",
    itemId: "item",
    field: "title",
    criterionIndex: -1,
    start,
    end,
    quote,
    body: "Feedback",
    batchId: "",
    createdAt: "now",
  };
}
test("highlights the exact selected UTF-16 range without changing plan text", () => {
  const text = "A 🎯 plan";
  const segments = annotatedSegments(text, [comment(2, 4, "🎯")]);
  expect(segments.map((segment) => segment.text).join("")).toBe(text);
  expect(
    segments.filter((segment) => segment.comments.length).map((segment) => segment.text),
  ).toEqual(["🎯"]);
});
test("overlapping comments create non-overlapping highlights with both comments", () => {
  const segments = annotatedSegments("abcdefgh", [
    comment(0, 5, "abcde", "first"),
    comment(3, 8, "defgh", "second"),
  ]);
  expect(segments.map((segment) => segment.text).join("")).toBe("abcdefgh");
  expect(
    segments.find((segment) => segment.text === "de")?.comments.map((item) => item.id),
  ).toEqual(["first", "second"]);
});
test("stale or out-of-bounds selections never highlight unrelated text", () => {
  const segments = annotatedSegments("New plan", [
    comment(0, 3, "Old"),
    comment(0, 100, "New plan"),
  ]);
  expect(segments).toEqual([{ start: 0, end: 8, text: "New plan", comments: [] }]);
});
