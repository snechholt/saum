import { useEffect, useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Factory as FactoryIcon, Plus, SquaresFour } from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import { deleteFactory, listFactories, type Factory } from "./api";
import { FactoryForm } from "./FactoryForm";
import { FactoryProjects } from "./FactoryProjects";

type FactoryList =
  | { status: "loading" }
  | { status: "ready"; factories: Factory[] }
  | { status: "failed" };
type Form = { mode: "create" } | { mode: "edit"; factory: Factory } | null;

export default function App() {
  const [list, setList] = useState<FactoryList>({ status: "loading" });
  const [selected, setSelected] = useState<Factory | null>(null);
  const [form, setForm] = useState<Form>(null);
  const [deleting, setDeleting] = useState<string | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    let ignore = false;
    listFactories()
      .then((factories) => {
        if (!ignore) setList({ status: "ready", factories });
      })
      .catch(() => {
        if (!ignore) setList({ status: "failed" });
      });
    return () => {
      ignore = true;
    };
  }, []);

  function factorySaved(factory: Factory) {
    setList((current) => {
      if (current.status !== "ready") return current;
      const exists = current.factories.some((item) => item.id === factory.id);
      return {
        status: "ready",
        factories: exists
          ? current.factories.map((item) => (item.id === factory.id ? factory : item))
          : [...current.factories, factory],
      };
    });
    setSelected((current) => (current?.id === factory.id ? factory : current));
    setForm(null);
  }

  async function removeFactory(factory: Factory) {
    if (deleting || !window.confirm(`Delete “${factory.name}”?`)) return;
    setDeleting(factory.id);
    setError("");
    try {
      await deleteFactory(factory.id);
      setList((current) =>
        current.status === "ready"
          ? {
              status: "ready",
              factories: current.factories.filter((item) => item.id !== factory.id),
            }
          : current,
      );
      setSelected((current) => (current?.id === factory.id ? null : current));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not delete factory.");
    } finally {
      setDeleting(null);
    }
  }

  const columns: Column<Factory>[] = [
    {
      id: "name",
      header: "Name",
      value: (factory) => factory.name,
      render: (factory) => (
        <button
          className="table-link"
          onClick={() => {
            setForm(null);
            setSelected(factory);
          }}
        >
          {factory.name}
        </button>
      ),
    },
  ];

  return (
    <div className="factory-layout">
      <aside className="sidebar" aria-label="Navigation">
        <div className="brand">
          <FactoryIcon size={23} weight="bold" />
          <strong>saum</strong>
        </div>
        <button
          className="nav-item"
          aria-current={selected ? undefined : "page"}
          onClick={() => {
            setSelected(null);
            setError("");
          }}
        >
          <SquaresFour size={17} />
          <span>Factories</span>
          {list.status === "ready" && <small>{list.factories.length}</small>}
        </button>
      </aside>
      <main className="workspace">
        {selected ? (
          <FactoryProjects key={selected.id} factory={selected} />
        ) : (
          <>
            <header className="page-title">
              <h1>Factories</h1>
              <Button
                variant="primary"
                size="sm"
                onClick={() => {
                  setForm({ mode: "create" });
                  setError("");
                }}
                disabled={list.status !== "ready" || form !== null || deleting !== null}
              >
                <Plus size={15} />
                New factory
              </Button>
            </header>
            {error && (
              <p className="error-message" role="alert">
                {error}
              </p>
            )}
            {form && (
              <FactoryForm
                key={form.mode === "edit" ? form.factory.id : "new"}
                {...(form.mode === "edit" ? { factory: form.factory } : {})}
                onSaved={factorySaved}
                onCancel={() => setForm(null)}
              />
            )}
            {list.status === "loading" ? (
              <output className="empty-state">Loading…</output>
            ) : list.status === "failed" ? (
              <p className="error-message" role="alert">
                Could not load factories. Reload to try again.
              </p>
            ) : (
              <DataTable
                rows={list.factories}
                columns={columns}
                rowKey={(factory) => factory.id}
                label="Factories"
                emptyMessage="No factories"
                actions={(factory) => (
                  <>
                    <Button
                      variant="ghost"
                      size="sm"
                      disabled={form !== null || deleting !== null}
                      aria-label={`Edit ${factory.name}`}
                      onClick={() => {
                        setForm({ mode: "edit", factory });
                        setError("");
                      }}
                    >
                      Edit
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      disabled={form !== null || deleting !== null}
                      aria-label={`Delete ${factory.name}`}
                      loading={deleting === factory.id}
                      onClick={() => {
                        void removeFactory(factory);
                      }}
                    >
                      Delete
                    </Button>
                  </>
                )}
              />
            )}
          </>
        )}
      </main>
    </div>
  );
}
