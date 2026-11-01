import { useEffect, useState } from "react";
import { Navigate, useParams } from "react-router";
import { getStory } from "./story-api";
import { sessionPath, storyPath } from "./paths";
// Existing deep links keep working, but there is only one factory-scoped workspace UI.
export function LegacyStoryRedirect() {
  const { storyId = "", sessionId } = useParams(),
    [target, setTarget] = useState(""),
    [error, setError] = useState("");
  useEffect(() => {
    let ignore = false;
    void getStory(storyId)
      .then((detail) => {
        if (!ignore)
          setTarget(
            sessionId
              ? sessionPath(detail.factory.id, storyId, sessionId)
              : storyPath(detail.factory.id, storyId),
          );
      })
      .catch((reason) => {
        if (!ignore) setError(reason instanceof Error ? reason.message : "Story not found.");
      });
    return () => {
      ignore = true;
    };
  }, [storyId, sessionId]);
  return target ? (
    <Navigate to={target} replace />
  ) : (
    <main className="workspace">
      {error ? <p role="alert">{error}</p> : <output>Loading…</output>}
    </main>
  );
}
