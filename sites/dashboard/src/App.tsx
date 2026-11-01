import { useEffect, useState, type FormEvent } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Badge } from "@cloudflare/kumo/components/badge";
import { CaretRight, Factory, FolderSimple, Plus, SquaresFour, X } from "@phosphor-icons/react";
import { createProject, listProjects, type Project } from "./api";

export default function App() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [saving, setSaving] = useState(false);
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    let ignore = false;
    listProjects()
      .then((items) => {
        if (!ignore) setProjects(items);
      })
      .catch(() => {
        if (!ignore) setError("Could not load projects. Make sure the Go server is running.");
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
      <aside className="sidebar" aria-label="Factory navigation">
        <div className="brand">
          <span className="brand-icon">
            <Factory size={21} weight="bold" />
          </span>
          <strong>saum</strong>
          <span className="local-badge">LOCAL</span>
        </div>
        <div className="factory-switch">
          <span className="factory-avatar">S</span>
          <div>
            <strong>Saum</strong>
            <small>Software factory</small>
          </div>
        </div>
        <p className="section-label">FACTORY</p>
        <div className="nav-item" aria-current="page">
          <SquaresFour size={17} />
          <span>Projects</span>
          <small>{projects.length}</small>
        </div>
        <div className="sidebar-footer">
          <span className="status-dot" />
          <span>Local factory</span>
        </div>
      </aside>

      <main className="workspace">
        <header className="workspace-header">
          <SquaresFour size={15} />
          <span>Saum</span>
          <CaretRight size={12} />
          <strong>Projects</strong>
        </header>
        <div className="workspace-content">
          <div className="eyebrow">
            <Factory size={14} /> SOFTWARE FACTORY
          </div>
          <div className="page-title">
            <div>
              <h1>Projects</h1>
              <p>Everything you’re building, in one place.</p>
            </div>
            <Button
              variant="primary"
              size="sm"
              onClick={() => setCreating(true)}
              disabled={loading || !!error || creating}
            >
              <Plus size={15} />
              New project
            </Button>
          </div>
          <div className="list-heading">
            <span>All projects</span>
            <Badge variant="neutral">{projects.length}</Badge>
          </div>
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
              <label htmlFor="project-name">Project name</label>
              <input
                id="project-name"
                value={name}
                onChange={(event) => setName(event.target.value)}
                maxLength={80}
                placeholder="e.g. ToDo App"
                disabled={saving}
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
                  Create project
                </Button>
              </div>
            </form>
          )}
          {loading ? (
            <output className="loading">Loading projects…</output>
          ) : projects.length === 0 ? (
            <section className="empty-state">
              <span className="empty-icon">
                <FolderSimple size={31} weight="duotone" />
              </span>
              <h2>A place for your next idea.</h2>
              <p>
                No projects yet. Start with a name.
                <br />
                The rest can come later.
              </p>
              {!creating && (
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => setCreating(true)}
                  disabled={!!error}
                >
                  <Plus size={14} />
                  Create your first project
                </Button>
              )}
            </section>
          ) : (
            <div className="project-grid">
              {projects.map((project) => (
                <article className="project-card" key={project.id}>
                  <span className="project-icon">
                    <FolderSimple size={22} weight="duotone" />
                  </span>
                  <h2>{project.name}</h2>
                  <p>Empty project</p>
                </article>
              ))}
            </div>
          )}
        </div>
      </main>

      <aside className="overview" aria-label="Factory overview">
        <header>Factory overview</header>
        <div className="overview-body">
          <span className="overview-icon">
            <Factory size={24} weight="duotone" />
          </span>
          <h2>Saum</h2>
          <p>Your software factory.</p>
          <div className="project-count">
            <span>Projects</span>
            <strong>{projects.length}</strong>
          </div>
          <p className="scope-note">
            Just projects, for now.
            <br />A project is an empty container for something you want to build.
          </p>
        </div>
      </aside>
    </div>
  );
}
