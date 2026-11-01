export function factoryPath(factoryId: string): string {
  return `/factories/${encodeURIComponent(factoryId)}`;
}
export function storyPath(factoryId: string, storyId: string): string {
  return `${factoryPath(factoryId)}/stories/${encodeURIComponent(storyId)}`;
}
export function sessionPath(factoryId: string, storyId: string, sessionId: string): string {
  return `${storyPath(factoryId, storyId)}/sessions/${encodeURIComponent(sessionId)}`;
}
export function auditPath(sessionId: string): string {
  return `/audit/${encodeURIComponent(sessionId)}`;
}
