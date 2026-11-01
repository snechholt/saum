import { expect, test } from "bun:test";
import { parseTranscript } from "../sites/dashboard/src/transcript";
import type { SessionEvent } from "../sites/dashboard/src/sessions";
function event(id: number, stream: string, data: unknown): SessionEvent {
  return {
    id,
    stream,
    data: Buffer.from(typeof data === "string" ? data : JSON.stringify(data)).toString("base64"),
  };
}

test("renders exposed Codex messages once as started items become completed", () => {
  const messages = [
    event(
      1,
      "stdout",
      JSON.stringify({
        type: "item.started",
        item: { id: "one", type: "agent_message", text: "" },
      }) + "\n",
    ),
    event(
      2,
      "stdout",
      JSON.stringify({
        type: "item.completed",
        item: { id: "one", type: "agent_message", text: "Hello 🎯" },
      }) + "\n",
    ),
  ];
  expect(parseTranscript(messages).entries).toEqual([
    { id: "item:one", kind: "message", text: "Hello 🎯" },
  ]);
});
test("handles UTF-8 split across persisted byte chunks", () => {
  const bytes = Buffer.from(
    JSON.stringify({
      type: "item.completed",
      item: { id: "one", type: "agent_message", text: "🎯" },
    }) + "\n",
  );
  const split = bytes.indexOf(Buffer.from("🎯")) + 1;
  const entries = parseTranscript([
    { id: 1, stream: "stdout", data: bytes.subarray(0, split).toString("base64") },
    { id: 2, stream: "stdout", data: bytes.subarray(split).toString("base64") },
  ]).entries;
  expect(entries).toEqual([{ id: "item:one", kind: "message", text: "🎯" }]);
});
test("shows codemode submissions and keeps diagnostics out of conversation", () => {
  const result = parseTranscript([
    event(1, "tool", {
      tool: "saum.submit_item",
      input: { key: "one", title: "Sign in" },
      result: { id: "saved" },
    }),
    event(2, "stderr", "diagnostic line"),
  ]);
  expect(result.entries).toHaveLength(1);
  expect(result.entries[0]?.kind).toBe("tool");
  expect(result.diagnostics).toBe("diagnostic line");
});
test("does not present a failed MCP call as a saved item", () => {
  const result = parseTranscript([
    event(
      1,
      "stdout",
      JSON.stringify({
        type: "item.completed",
        item: {
          id: "tool",
          type: "mcp_tool_call",
          tool: "execute",
          arguments: { code: "bad()" },
          result: { isError: true },
          status: "completed",
        },
      }) + "\n",
    ),
  ]);
  expect(result.entries[0]).toMatchObject({
    kind: "tool",
    status: "failed",
    error: "Tool execution failed",
  });
});
