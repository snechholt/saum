import { describe, expect, test } from "bun:test";
import { tableRows } from "./data";

const rows = [
  { name: "Beta", count: 10 },
  { name: "alpha", count: 2 },
  { name: "Alpha 2", count: 1 },
];
const columns = [
  { id: "name", value: (row: (typeof rows)[number]) => row.name },
  { id: "count", value: (row: (typeof rows)[number]) => row.count },
];

describe("table rows", () => {
  test("filters across columns, ignoring case and surrounding whitespace", () => {
    expect(tableRows(rows, columns, " ALPHA ", null).map((row) => row.name)).toEqual([
      "alpha",
      "Alpha 2",
    ]);
    expect(tableRows(rows, columns, "10", null)).toEqual([{ name: "Beta", count: 10 }]);
    expect(tableRows(rows, columns, "missing", null)).toEqual([]);
  });
  test("sorts names in either direction without mutating the input", () => {
    expect(
      tableRows(rows, columns, "", { column: "name", direction: "ascending" }).map(
        (row) => row.name,
      ),
    ).toEqual(["alpha", "Alpha 2", "Beta"]);
    expect(
      tableRows(rows, columns, "", { column: "name", direction: "descending" }).map(
        (row) => row.name,
      ),
    ).toEqual(["Beta", "Alpha 2", "alpha"]);
    expect(rows[0]?.name).toBe("Beta");
  });
  test("sorts numeric columns numerically", () => {
    expect(
      tableRows(rows, columns, "", { column: "count", direction: "ascending" }).map(
        (row) => row.count,
      ),
    ).toEqual([1, 2, 10]);
  });
  test("retains source order without a valid sort column", () => {
    expect(tableRows(rows, columns, "", null)).toEqual(rows);
    expect(tableRows(rows, columns, "", { column: "missing", direction: "ascending" })).toEqual(
      rows,
    );
  });
});
