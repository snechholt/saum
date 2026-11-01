import { useEffect, useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Sidebar } from "@cloudflare/kumo/components/sidebar";
import {
  Factory as FactoryIcon,
  PencilSimple,
  Plus,
  SquaresFour,
  Trash,
} from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import { deleteFactory, listFactories, type Factory } from "./api";
import { FactoryForm } from "./FactoryForm";
import { FactoryStories } from "./FactoryStories";
import { DeleteFactoryDialog } from "./DeleteFactoryDialog";

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
  const [pendingDeletion, setPendingDeletion] = useState<Factory | null>(null);

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
    if (deleting) return;
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
      setPendingDeletion(null);
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
    <Sidebar.Provider
      defaultOpen
      defaultWidth={224}
      mobileBreakpoint={640}
      collapsible="offcanvas"
      className="factory-layout"
    >
      <Sidebar aria-label="Navigation">
        <Sidebar.Header>
          <div className="brand">
            <FactoryIcon size={22} />
            <span>Saum</span>
          </div>
        </Sidebar.Header>
        <Sidebar.Content>
          <Sidebar.Menu>
            <Sidebar.MenuButton
              className="nav-item"
              icon={SquaresFour}
              active={!selected}
              aria-current={selected ? undefined : "page"}
              onClick={() => {
                setSelected(null);
                setError("");
              }}
            >
              Factories
              {list.status === "ready" && (
                <Sidebar.MenuBadge>{list.factories.length}</Sidebar.MenuBadge>
              )}
            </Sidebar.MenuButton>
          </Sidebar.Menu>
        </Sidebar.Content>
      </Sidebar>
      <main className="workspace">
        <div className="mobile-navigation">
          <Sidebar.Trigger />
          <span>Saum</span>
        </div>
        {selected ? (
          <FactoryStories key={selected.id} factory={selected} />
        ) : (
          <>
            <header className="page-title">
              <h1>Factories</h1>
              <Button
                variant="primary"
                size="base"
                shape="square"
                aria-label="New factory"
                title="New factory"
                onClick={() => {
                  setForm({ mode: "create" });
                  setError("");
                }}
                disabled={list.status !== "ready" || form !== null || deleting !== null}
              >
                <Plus size={15} />
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
                      shape="square"
                      title="Edit"
                      aria-label={`Edit ${factory.name}`}
                      onClick={() => {
                        setForm({ mode: "edit", factory });
                        setError("");
                      }}
                    >
                      <PencilSimple size={15} />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      disabled={form !== null || deleting !== null}
                      shape="square"
                      title="Delete"
                      aria-label={`Delete ${factory.name}`}
                      loading={deleting === factory.id}
                      onClick={() => {
                        setPendingDeletion(factory);
                      }}
                    >
                      <Trash size={15} />
                    </Button>
                  </>
                )}
              />
            )}
          </>
        )}
      </main>
      <DeleteFactoryDialog
        factory={pendingDeletion}
        busy={deleting !== null}
        onCancel={() => setPendingDeletion(null)}
        onConfirm={(factory) => {
          void removeFactory(factory);
        }}
      />
    </Sidebar.Provider>
  );
}
