import type { SessionEvent } from "./sessions";

export type TranscriptEntry =
  | { id: string; kind: "message"; text: string }
  | { id: string; kind: "command"; command: string; output: string; status: string }
  | {
      id: string;
      kind: "tool";
      name: string;
      input: Record<string, unknown>;
      result: unknown;
      error: string;
      status: string;
    };

function object(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}
const text = (value: unknown) => (typeof value === "string" ? value : "");
function json(value: string): unknown {
  try {
    const parsed: unknown = JSON.parse(value);
    return parsed;
  } catch {
    return null;
  }
}
function functionName(value: string) {
  return value === "saum_item_submit"
    ? "saum.submit_item"
    : value === "saum_plan_complete"
      ? "saum.complete_plan"
      : value;
}
export function parseTranscript(events: readonly SessionEvent[]): {
  entries: TranscriptEntry[];
  diagnostics: string;
} {
  const entries = new Map<string, TranscriptEntry>();
  const decoders = new Map<string, TextDecoder>();
  let stdout = "",
    diagnostics = "";
  function event(value: unknown) {
    if (!object(value) || !object(value["item"])) return;
    const item = value["item"],
      type = text(item["type"]),
      id = `item:${text(item["id"])}`;
    if (type === "agent_message") {
      const content = text(item["text"]);
      if (content) entries.set(id, { id, kind: "message", text: content });
    } else if (type === "command_execution") {
      entries.set(id, {
        id,
        kind: "command",
        command: text(item["command"]),
        output: text(item["aggregated_output"]),
        status: text(item["status"]),
      });
    } else if (type === "mcp_tool_call") {
      const input = object(item["arguments"]) ? item["arguments"] : {};
      const name = functionName(text(item["tool"]));
      const toolID = name === "saum.submit_item" ? `function:${text(input["key"])}` : id;
      const rawResult = item["result"];
      const resultFailed = object(rawResult) && rawResult["isError"] === true;
      const error =
        (object(item["error"]) ? text(item["error"]["message"]) : text(item["error"])) ||
        (resultFailed ? "Tool execution failed" : "");
      entries.set(toolID, {
        id: toolID,
        kind: "tool",
        name,
        input,
        result: item["result"],
        error,
        status: error ? "failed" : text(item["status"]),
      });
    }
  }
  for (const row of events) {
    let decoder = decoders.get(row.stream);
    if (!decoder) {
      decoder = new TextDecoder();
      decoders.set(row.stream, decoder);
    }
    const chunk = decoder.decode(
      Uint8Array.from(atob(row.data), (character) => character.charCodeAt(0)),
      { stream: true },
    );
    if (row.stream === "stderr") {
      diagnostics += chunk;
      continue;
    }
    if (row.stream === "stdout") {
      stdout += chunk;
      let boundary = stdout.indexOf("\n");
      while (boundary >= 0) {
        event(json(stdout.slice(0, boundary)));
        stdout = stdout.slice(boundary + 1);
        boundary = stdout.indexOf("\n");
      }
    } else if (row.stream === "tool") {
      const value = json(chunk);
      if (!object(value)) continue;
      const name = functionName(text(value["tool"]));
      if (
        ![
          "saum.submit_item",
          "saum.update_item",
          "saum.remove_item",
          "saum.resolve_feedback",
          "saum.complete_plan",
          // Decode historical events without re-exposing these old functions to agents.
          "saum.start_item",
          "saum.complete_item",
          "saum.block_item",
          "saum.complete_implementation",
          "saum.work_item.submit_acceptance_criteria",
          "saum.work_item.request_review",
          "saum.work_item.report_blocker",
          "saum.work_item.submit_executive_summary",
        ].includes(name)
      )
        continue;
      const input = object(value["input"]) ? value["input"] : {};
      const id =
        name === "saum.submit_item"
          ? `function:${text(input["key"])}`
          : name === "saum.complete_plan"
            ? `function:${name}`
            : `function:${name}:${row.id}`;
      entries.set(id, {
        id,
        kind: "tool",
        name,
        input,
        result: value["result"],
        error: text(value["error"]),
        status: value["error"] ? "failed" : "completed",
      });
    }
  }
  if (stdout.trim()) event(json(stdout));
  return { entries: [...entries.values()], diagnostics };
}
