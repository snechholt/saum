import { useEffect, useState } from "react";
import { DataTable, type Column } from "@saum/ui";
import { listProjects, type Factory, type Project } from "./api";

const columns: Column<Project>[] = [
  { id: "name", header: "Name", value: (project) => project.name },
];

type ProjectList =
  | { status: "loading" }
  | { status: "ready"; projects: Project[] }
  | { status: "failed" };

export function FactoryProjects({ factory }: { factory: Factory }) {
  const [list, setList] = useState<ProjectList>({ status: "loading" });
  useEffect(() => {
    let ignore = false;
    listProjects(factory.id)
      .then((projects) => {
        if (!ignore) setList({ status: "ready", projects });
      })
      .catch(() => {
        if (!ignore) setList({ status: "failed" });
      });
    return () => {
      ignore = true;
    };
  }, [factory.id]);

  return (
    <>
      <header className="page-title">
        <h1>{factory.name}</h1>
      </header>
      <h2 className="section-title">Projects</h2>
      {list.status === "loading" ? (
        <output className="empty-state">Loading…</output>
      ) : list.status === "failed" ? (
        <p className="error-message" role="alert">
          Could not load projects.
        </p>
      ) : (
        <DataTable
          rows={list.projects}
          columns={columns}
          rowKey={(project) => project.id}
          label="Projects"
          emptyMessage="No projects"
        />
      )}
    </>
  );
}
