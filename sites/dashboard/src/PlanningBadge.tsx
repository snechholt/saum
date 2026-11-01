import { Badge } from "@cloudflare/kumo/components/badge";
import type { PlanningStatus } from "./api";
import { planningLabel } from "./sessions";

export function PlanningBadge({ status }: { status: PlanningStatus }) {
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
      {planningLabel(status)}
    </Badge>
  );
}
