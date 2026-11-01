import { useCallback, useEffect, useState, type FormEvent } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Factory, FolderSimple, Plus, SquaresFour, X } from "@phosphor-icons/react";
import { createProject, listProjects, type Project } from "./api";

export default function App() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [saving, setSaving] = useState(false);
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [loadFailed, setLoadFailed] = useState(false);
  const focusName = useCallback((element: HTMLInputElement | null) => element?.focus(), []);

  useEffect(() => {
    let ignore = false;
    listProjects()
      .then((items) => {
        if (!ignore) setProjects(items);
      })
      .catch(() => {
        if (!ignore) {
          setLoadFailed(true);
          setError("Could not load projects. Reload to try again.");
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!name.trim() || saving) return;
    setSaving(true);
    setError("");
    try {
      const project = await createProject(name);
      setProjects((current) => [...current, project]);
      setName("");
      setCreating(false);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not create project.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="factory-layout">
      <aside className="sidebar" aria-label="Navigation">
        <div className="brand">
          <Factory size={23} weight="bold" />
          <strong>saum</strong>
        </div>
        <div className="nav-item" aria-current="page">
          <SquaresFour size={17} />
          <span>Projects</span>
          <small>{projects.length}</small>
        </div>
      </aside>
      <main className="workspace">
        <header className="page-title">
          <h1>Projects</h1>
          <Button
            variant="primary"
            size="sm"
            onClick={() => setCreating(true)}
            disabled={loading || loadFailed || creating}
          >
            <Plus size={15} />
            New project
          </Button>
        </header>
        {error && (
          <div className="error-message" role="alert">
            {error}
            <button onClick={() => setError("")} aria-label="Dismiss error">
              <X size={14} />
            </button>
          </div>
        )}
        {creating && (
          <form
            className="create-project"
            onSubmit={(event) => {
              void submit(event);
            }}
          >
            <label htmlFor="project-name">Name</label>
            <input
              id="project-name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              maxLength={80}
              disabled={saving}
              ref={focusName}
            />
            <div>
              <Button
                variant="ghost"
                size="sm"
                type="button"
                onClick={() => {
                  setCreating(false);
                  setName("");
                }}
                disabled={saving}
              >
                Cancel
              </Button>
              <Button
                variant="primary"
                size="sm"
                type="submit"
                loading={saving}
                disabled={!name.trim() || saving}
              >
                Create
              </Button>
            </div>
          </form>
        )}
        {loading ? (
          <output className="empty-state">Loading…</output>
        ) : loadFailed ? null : projects.length === 0 ? (
          <p className="empty-state">No projects</p>
        ) : (
          <div className="project-grid">
            {projects.map((project) => (
              <article className="project-card" key={project.id}>
                <FolderSimple size={22} weight="duotone" />
                <h2>{project.name}</h2>
              </article>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
