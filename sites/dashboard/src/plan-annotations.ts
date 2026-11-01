import type { PlanComment } from "./sessions";
export type AnnotatedSegment = {
  start: number;
  end: number;
  text: string;
  comments: PlanComment[];
};
export function annotatedSegments(
  text: string,
  comments: readonly PlanComment[],
): AnnotatedSegment[] {
  const valid = comments.filter(
    (comment) =>
      comment.start >= 0 &&
      comment.end <= text.length &&
      comment.start < comment.end &&
      text.slice(comment.start, comment.end) === comment.quote,
  );
  const boundaries = [
    ...new Set([0, text.length, ...valid.flatMap((comment) => [comment.start, comment.end])]),
  ].sort((a, b) => a - b);
  return boundaries.slice(0, -1).map((start, index) => {
    const end = boundaries[index + 1] ?? text.length;
    return {
      start,
      end,
      text: text.slice(start, end),
      comments: valid.filter((comment) => comment.start < end && comment.end > start),
    };
  });
}
