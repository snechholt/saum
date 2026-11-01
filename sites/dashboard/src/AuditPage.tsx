import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "@cloudflare/kumo/components/button";
import { Select } from "@cloudflare/kumo/components/select";
import { ArrowClockwise, CaretDown, Terminal } from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import type { Factory, SessionKind } from "./api";
import { getAudit, type AuditEntry, type AuditPage as AuditResult } from "./audit-api";
import { PlanningBadge } from "./PlanningBadge";
import { auditPath, factoryPath } from "./paths";
import { AuditTranscript } from "./AuditTranscript";
const roles = { planning: "Planner", implementation: "Goal agent", summary: "Factory agent" };
const columns: Column<AuditEntry>[] = [
  {
    id: "factory",
    header: "Factory",
    value: (entry) => entry.factory.name,
    render: (entry) => (
      <Link className="table-link" to={factoryPath(entry.factory.id)}>
        {entry.factory.name}
      </Link>
    ),
  },
  { id: "goal", header: "Goal", value: (entry) => entry.storyTitle },
  { id: "workItem", header: "Work item", value: (entry) => entry.workItemTitle || "—" },
  { id: "role", header: "Agent", value: (entry) => roles[entry.session.kind] },
  {
    id: "status",
    header: "Status",
    value: (entry) =>
      entry.session.error.toLowerCase().startsWith("blocked:") ? "blocked" : entry.session.status,
    render: (entry) => (
      <PlanningBadge
        status={entry.session.status}
        kind={entry.session.kind}
        goalStatus={
          entry.session.error.toLowerCase().startsWith("blocked:") ? "blocked" : undefined
        }
      />
    ),
  },
  {
    id: "thread",
    header: "Thread",
    value: (entry) => entry.session.threadId,
    render: (entry) => (
      <code className="workspace-path" title={entry.session.threadId}>
        {entry.session.threadId
          ? `${entry.session.threadId.slice(0, 8)}…${entry.session.threadId.slice(-6)}`
          : "—"}
      </code>
    ),
  },
  {
    id: "started",
    header: "Started",
    value: (entry) => Date.parse(entry.session.createdAt),
    render: (entry) => (
      <time dateTime={entry.session.createdAt}>
        {new Date(entry.session.createdAt).toLocaleString()}
      </time>
    ),
  },
];
function AuditList({ factories }: { factories: Factory[] }) {
  const navigate = useNavigate(),
    [factory, setFactory] = useState("all"),
    [role, setRole] = useState<SessionKind | "all">("all"),
    [page, setPage] = useState<AuditResult | null>(null),
    [error, setError] = useState(""),
    [loading, setLoading] = useState(true),
    [version, setVersion] = useState(0),
    generation = useRef(0);
  useEffect(() => {
    const ticket = ++generation.current;
    void getAudit(factory, role)
      .then((result) => {
        if (ticket === generation.current) setPage(result);
      })
      .catch((reason) => {
        if (ticket === generation.current)
          setError(reason instanceof Error ? reason.message : "Could not load audit.");
      })
      .finally(() => {
        if (ticket === generation.current) setLoading(false);
      });
    return () => {
      generation.current = ticket + 1;
    };
  }, [factory, role, version]);
  function reset() {
    generation.current++;
    setPage(null);
    setError("");
    setLoading(true);
  }
  async function more() {
    if (!page?.nextCursor || loading) return;
    const ticket = generation.current;
    setError("");
    setLoading(true);
    try {
      const result = await getAudit(factory, role, page.nextCursor);
      if (ticket === generation.current)
        setPage({ entries: [...page.entries, ...result.entries], nextCursor: result.nextCursor });
    } catch (reason) {
      if (ticket === generation.current)
        setError(reason instanceof Error ? reason.message : "Could not load more sessions.");
    } finally {
      if (ticket === generation.current) setLoading(false);
    }
  }
  return (
    <>
      <header className="page-title">
        <h1>Audit</h1>
        <Button
          variant="secondary"
          shape="square"
          aria-label="Refresh audit"
          title="Refresh audit"
          disabled={loading}
          onClick={() => {
            reset();
            setVersion((current) => current + 1);
          }}
        >
          <ArrowClockwise size={16} />
        </Button>
      </header>
      <div className="audit-filters">
        <Select<string>
          aria-label="Filter by factory"
          value={factory}
          items={[
            { value: "all", label: "All factories" },
            ...factories.map((factory) => ({ value: factory.id, label: factory.name })),
          ]}
          onValueChange={(value) => {
            if (value && value !== factory) {
              reset();
              setFactory(value);
            }
          }}
        />
        <Select<SessionKind | "all">
          aria-label="Filter by agent"
          value={role}
          items={[
            { value: "all", label: "All agents" },
            { value: "planning", label: "Planner" },
            { value: "implementation", label: "Goal agent" },
            { value: "summary", label: "Factory agent" },
          ]}
          onValueChange={(value) => {
            if (value && value !== role) {
              reset();
              setRole(value);
            }
          }}
        />
      </div>
      {error && (
        <p className="error-message" role="alert">
          {error}
        </p>
      )}
      {page ? (
        <DataTable
          rows={page.entries}
          columns={columns}
          rowKey={(entry) => entry.session.id}
          label="Agent sessions"
          emptyMessage="No agent sessions"
          actions={(entry) => (
            <Button
              variant="ghost"
              size="sm"
              shape="square"
              aria-label={`Open transcript ${entry.session.id}`}
              title="Open transcript"
              onClick={() => {
                void navigate(auditPath(entry.session.id));
              }}
            >
              <Terminal size={16} />
            </Button>
          )}
        />
      ) : loading ? (
        <output>Loading…</output>
      ) : null}
      {!!page?.nextCursor && (
        <div className="form-actions">
          <Button
            variant="secondary"
            shape="square"
            aria-label="Load more sessions"
            title="Load more sessions"
            disabled={loading}
            loading={loading}
            onClick={() => {
              void more();
            }}
          >
            <CaretDown size={16} />
          </Button>
        </div>
      )}
    </>
  );
}
export default function AuditPage({
  factories,
  sessionId,
}: {
  factories: Factory[];
  sessionId: string;
}) {
  return sessionId ? (
    <AuditTranscript key={sessionId} id={sessionId} />
  ) : (
    <AuditList factories={factories} />
  );
}
