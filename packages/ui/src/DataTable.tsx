import { useState, type ReactNode } from "react";
import { Input } from "@cloudflare/kumo/components/input";
import { LayerCard } from "@cloudflare/kumo/components/layer-card";
import { Table } from "@cloudflare/kumo/components/table";
import { ArrowDown, ArrowUp, ArrowsDownUp } from "@phosphor-icons/react";
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
      <div className="table-toolbar">
        <Input
          className="table-filter"
          type="search"
          aria-label={`Filter ${label.toLowerCase()}`}
          placeholder="Filter…"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
        />
      </div>
      <LayerCard className="table-surface">
        <div className="table-scroll">
          <Table aria-label={label}>
            <Table.Header variant="compact">
              <Table.Row>
                {columns.map((column) => (
                  <Table.Head
                    key={column.id}
                    aria-sort={sort?.column === column.id ? sort.direction : "none"}
                  >
                    <button
                      className="table-sort"
                      type="button"
                      onClick={() => sortColumn(column.id)}
                      aria-label={`Sort by ${column.header.toLowerCase()}`}
                    >
                      {column.header}
                      {sort?.column === column.id ? (
                        sort.direction === "ascending" ? (
                          <ArrowUp size={14} aria-hidden />
                        ) : (
                          <ArrowDown size={14} aria-hidden />
                        )
                      ) : (
                        <ArrowsDownUp size={14} aria-hidden />
                      )}
                    </button>
                  </Table.Head>
                ))}
                {actions && (
                  <Table.Head>
                    <span className="table-sr-only">Actions</span>
                  </Table.Head>
                )}
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {visible.length === 0 ? (
                <Table.Row>
                  <Table.Cell className="table-empty" colSpan={columns.length + (actions ? 1 : 0)}>
                    {rows.length === 0 ? emptyMessage : "No matches"}
                  </Table.Cell>
                </Table.Row>
              ) : (
                visible.map((row) => (
                  <Table.Row key={rowKey(row)}>
                    {columns.map((column) => (
                      <Table.Cell key={column.id}>
                        {column.render ? column.render(row) : column.value(row)}
                      </Table.Cell>
                    ))}
                    {actions && (
                      <Table.Cell className="table-actions">
                        <div className="table-row-actions">{actions(row)}</div>
                      </Table.Cell>
                    )}
                  </Table.Row>
                ))
              )}
            </Table.Body>
          </Table>
        </div>
      </LayerCard>
    </div>
  );
}
