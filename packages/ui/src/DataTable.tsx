import { useState, type ReactNode } from "react";
import { tableRows, type Sort, type ValueColumn } from "./data";

export type Column<T> = ValueColumn<T> & {
  header: string;
  render?: (row: T) => ReactNode;
};

type Props<T> = {
  rows: readonly T[];
  columns: readonly Column<T>[];
  rowKey: (row: T) => string;
  label: string;
  emptyMessage: string;
  actions?: (row: T) => ReactNode;
};

export function DataTable<T>({ rows, columns, rowKey, label, emptyMessage, actions }: Props<T>) {
  const [query, setQuery] = useState("");
  const [sort, setSort] = useState<Sort | null>(null);
  const visible = tableRows(rows, columns, query, sort);

  function sortColumn(column: string) {
    setSort((current) => ({
      column,
      direction:
        current?.column === column && current.direction === "ascending"
          ? "descending"
          : "ascending",
    }));
  }

  return (
    <div className="data-table">
      <input
        className="table-filter"
        type="search"
        aria-label={`Filter ${label.toLowerCase()}`}
        placeholder="Filter…"
        value={query}
        onChange={(event) => setQuery(event.target.value)}
      />
      <div className="table-scroll">
        <table aria-label={label}>
          <thead>
            <tr>
              {columns.map((column) => (
                <th
                  key={column.id}
                  scope="col"
                  aria-sort={sort?.column === column.id ? sort.direction : "none"}
                >
                  <button
                    type="button"
                    onClick={() => sortColumn(column.id)}
                    aria-label={`Sort by ${column.header.toLowerCase()}`}
                  >
                    {column.header}
                    <span aria-hidden="true">
                      {sort?.column === column.id
                        ? sort.direction === "ascending"
                          ? "↑"
                          : "↓"
                        : "↕"}
                    </span>
                  </button>
                </th>
              ))}
              {actions && (
                <th scope="col">
                  <span className="table-sr-only">Actions</span>
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {visible.length === 0 ? (
              <tr>
                <td className="table-empty" colSpan={columns.length + (actions ? 1 : 0)}>
                  {rows.length === 0 ? emptyMessage : "No matches"}
                </td>
              </tr>
            ) : (
              visible.map((row) => (
                <tr key={rowKey(row)}>
                  {columns.map((column) => (
                    <td key={column.id}>
                      {column.render ? column.render(row) : column.value(row)}
                    </td>
                  ))}
                  {actions && <td className="table-actions">{actions(row)}</td>}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
