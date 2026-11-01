import { Badge } from "@cloudflare/kumo/components/badge";
import type { GoalStatus, PlanningStatus, SessionKind } from "./api";
import { planningLabel } from "./sessions";
export function PlanningBadge({
  status,
  kind = "planning",
  goalStatus,
}: {
  status: PlanningStatus;
  kind?: SessionKind;
  goalStatus?: GoalStatus | undefined;
}) {
  const label =
    goalStatus === "completed"
      ? "Completed"
      : goalStatus === "awaiting_review"
        ? "Needs review"
        : goalStatus === "blocked"
          ? "Blocked"
          : goalStatus === "failed"
            ? "Failed"
            : planningLabel(status, kind);
  const success = goalStatus === "completed" || (kind === "planning" && status === "completed");
  const error =
    goalStatus === "blocked" ||
    goalStatus === "failed" ||
    status === "failed" ||
    status === "interrupted";
  const info =
    status === "running" ||
    status === "preparing" ||
    (kind !== "planning" && status === "completed");
  return (
    <Badge variant={error ? "error" : success ? "success" : info ? "info" : "neutral"}>
      {label}
    </Badge>
  );
}
