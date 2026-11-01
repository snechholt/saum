import { Badge } from "@cloudflare/kumo/components/badge";
import { DataTable, type Column } from "@saum/ui";
import type { Implementation, ImplementationItem, Session } from "./sessions";
import { PlanningBadge } from "./PlanningBadge";
const labels = { pending: "Queued", running: "Building", completed: "Done", blocked: "Blocked" };
const columns: Column<ImplementationItem>[] = [
  { id: "title", header: "Item", value: (item) => item.title },
  {
    id: "status",
    header: "Status",
    value: (item) => item.status,
    render: (item) => (
      <Badge
        variant={
          item.status === "completed"
            ? "success"
            : item.status === "blocked"
              ? "error"
              : item.status === "running"
                ? "info"
                : "neutral"
        }
      >
        {labels[item.status]}
      </Badge>
    ),
  },
  {
    id: "summary",
    header: "Result",
    value: (item) => item.summary,
    render: (item) => <span className="implementation-result">{item.summary || "—"}</span>,
  },
];
export function ImplementationProgress({
  session,
  implementation,
}: {
  session: Session;
  implementation: Implementation;
}) {
  const blocked = implementation.items.some((item) => item.status === "blocked");
  return (
    <section aria-label="Implementation" className="implementation-progress">
      <header className="section-toolbar">
        <h2 className="section-title">Implementation</h2>
        <div className="session-actions">
          <span>
            {implementation.items.filter((item) => item.status === "completed").length} /{" "}
            {implementation.items.length}
          </span>
          {blocked ? (
            <Badge variant="error">Blocked</Badge>
          ) : (
            <PlanningBadge status={session.status} kind="implementation" />
          )}
        </div>
      </header>
      <DataTable
        rows={implementation.items}
        columns={columns}
        rowKey={(item) => item.key}
        label="Implementation items"
        emptyMessage="No approved items"
      />
      {implementation.summary && <p className="implementation-result">{implementation.summary}</p>}
    </section>
  );
}
