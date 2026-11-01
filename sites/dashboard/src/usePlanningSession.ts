import { useCallback, useEffect, useRef, useState } from "react";
import { getSession, getSessionEvents, type SessionDetail, type SessionEvent } from "./sessions";
type Snapshot = { id: string; detail: SessionDetail | null; events: SessionEvent[]; error: string };
export function usePlanningSession(id: string) {
  const [snapshot, setSnapshot] = useState<Snapshot>({ id, detail: null, events: [], error: "" });
  const generation = useRef(0);
  const refresh = useCallback(async () => {
    const ticket = ++generation.current;
    const detail = await getSession(id);
    if (ticket !== generation.current) return;
    setSnapshot((current) => ({
      id,
      detail,
      events: current.id === id ? current.events : [],
      error: "",
    }));
  }, [id]);
  useEffect(() => {
    if (!id) return undefined;
    let ignore = false,
      after = 0;
    let timer: ReturnType<typeof setTimeout>;
    async function poll() {
      const ticket = generation.current;
      try {
        const [detail, batch] = await Promise.all([getSession(id), getSessionEvents(id, after)]);
        if (ignore) return;
        after = batch.at(-1)?.id ?? after;
        setSnapshot((current) => ({
          id,
          detail: ticket === generation.current ? detail : current.detail,
          events: [...(current.id === id ? current.events : []), ...batch],
          error: "",
        }));
        timer = setTimeout(
          () => {
            void poll();
          },
          batch.length === 200 ? 0 : 1000,
        );
      } catch (reason) {
        if (!ignore) {
          setSnapshot((current) => ({
            ...(current.id === id ? current : { id, detail: null, events: [] }),
            error: reason instanceof Error ? reason.message : "Could not load session.",
          }));
          timer = setTimeout(() => {
            void poll();
          }, 2000);
        }
      }
    }
    void poll();
    return () => {
      ignore = true;
      clearTimeout(timer);
    };
  }, [id]);
  return {
    ...(snapshot.id === id ? snapshot : { id, detail: null, events: [], error: "" }),
    refresh,
  };
}
