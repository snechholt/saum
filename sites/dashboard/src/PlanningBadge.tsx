import { Badge } from "@cloudflare/kumo/components/badge";
import type { PlanningStatus, SessionKind } from "./api";
import { planningLabel } from "./sessions";

export function PlanningBadge({
  status,
  kind = "planning",
}: {
  status: PlanningStatus;
  kind?: SessionKind;
}) {
  return (
    <Badge
      variant={
        status === "completed"
          ? "success"
          : status === "failed" || status === "interrupted"
            ? "error"
            : status === "running" || status === "preparing"
              ? "info"
              : "neutral"
      }
    >
      {planningLabel(status, kind)}
    </Badge>
  );
}
