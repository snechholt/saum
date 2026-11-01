import { useEffect, useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Factory as FactoryIcon, Plus, SquaresFour } from "@phosphor-icons/react";
import { listFactories, type Factory } from "./api";
import { NewFactoryForm } from "./NewFactoryForm";
import { FactoryProjects } from "./FactoryProjects";

type FactoryList =
  | { status: "loading" }
  | { status: "ready"; factories: Factory[] }
  | { status: "failed" };

export default function App() {
  const [list, setList] = useState<FactoryList>({ status: "loading" });
  const [selected, setSelected] = useState<Factory | null>(null);
  const [creating, setCreating] = useState(false);

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

  function factoryCreated(factory: Factory) {
    setList((current) =>
      current.status === "ready"
        ? { status: "ready", factories: [...current.factories, factory] }
        : current,
    );
    setCreating(false);
  }

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
          onClick={() => setSelected(null)}
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
                onClick={() => setCreating(true)}
                disabled={list.status !== "ready" || creating}
              >
                <Plus size={15} />
                New factory
              </Button>
            </header>
            {creating && (
              <NewFactoryForm onCreated={factoryCreated} onCancel={() => setCreating(false)} />
            )}
            {list.status === "loading" ? (
              <output className="empty-state">Loading…</output>
            ) : list.status === "failed" ? (
              <p className="error-message" role="alert">
                Could not load factories. Reload to try again.
              </p>
            ) : list.factories.length === 0 ? (
              <p className="empty-state">No factories</p>
            ) : (
              <div className="item-grid">
                {list.factories.map((factory) => (
                  <button
                    className="item-card factory-card"
                    key={factory.id}
                    onClick={() => {
                      setCreating(false);
                      setSelected(factory);
                    }}
                  >
                    <FactoryIcon size={22} weight="duotone" />
                    <span className="item-name">{factory.name}</span>
                  </button>
                ))}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
