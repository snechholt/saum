export type Sort = { column: string; direction: "ascending" | "descending" };
export type ValueColumn<T> = {
  id: string;
  value: (row: T) => string | number;
};
const collator = new Intl.Collator(undefined, { numeric: true, sensitivity: "base" });

export function tableRows<T>(
  rows: readonly T[],
  columns: readonly ValueColumn<T>[],
  query: string,
  sort: Sort | null,
): T[] {
  const search = query.trim().toLocaleLowerCase();
  const filtered = rows.filter(
    (row) =>
      !search ||
      columns.some((column) => String(column.value(row)).toLocaleLowerCase().includes(search)),
  );
  const column = columns.find((column) => column.id === sort?.column);
  if (!sort || !column) return filtered;
  const direction = sort.direction === "ascending" ? 1 : -1;
  return filtered.sort((left, right) => {
    const a = column.value(left);
    const b = column.value(right);
    const order =
      typeof a === "number" && typeof b === "number"
        ? a - b
        : collator.compare(String(a), String(b));
    return direction * order;
  });
}
