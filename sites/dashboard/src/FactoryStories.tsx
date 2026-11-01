import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "@cloudflare/kumo/components/button";
import { ArrowClockwise, Plus, Terminal } from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import { listStories, type Factory, type UserStory } from "./api";
import { isActive, retryPlanning, retryImplementation } from "./sessions";
import { PlanningBadge } from "./PlanningBadge";
import { UserStoryForm } from "./UserStoryForm";

const columns: Column<UserStory>[] = [
  {
    id: "title",
    header: "Title",
    value: (story) => story.title,
    render: (story) => (
      <Link className="table-link story-title" to={`/stories/${encodeURIComponent(story.id)}`}>
        {story.title}
      </Link>
    ),
  },
  {
    id: "goal",
    header: "Goal",
    value: (story) => story.goal,
    render: (story) => <span className="story-goal">{story.goal || "—"}</span>,
  },
  {
    id: "planning",
    header: "Status",
    value: (story) => story.planningStatus,
    render: (story) => <PlanningBadge status={story.planningStatus} kind={story.sessionKind} />,
  },
];
type StoryList =
  | { status: "loading" }
  | { status: "ready"; stories: UserStory[] }
  | { status: "failed" };

export function FactoryStories({ factory }: { factory: Factory }) {
  const [list, setList] = useState<StoryList>({ status: "loading" });
  const [creating, setCreating] = useState(false);
  const navigate = useNavigate();
  const [retrying, setRetrying] = useState("");
  const [error, setError] = useState("");
  useEffect(() => {
    let ignore = false;
    let timer: ReturnType<typeof setTimeout>;
    async function poll() {
      try {
        const stories = await listStories(factory.id);
        if (!ignore) setList({ status: "ready", stories });
      } catch {
        if (!ignore)
          setList((current) => (current.status === "loading" ? { status: "failed" } : current));
      }
      if (!ignore)
        timer = setTimeout(() => {
          void poll();
        }, 1200);
    }
    void poll();
    return () => {
      ignore = true;
      clearTimeout(timer);
    };
  }, [factory.id]);

  function storyCreated(story: UserStory) {
    setList((current) =>
      current.status === "ready"
        ? { status: "ready", stories: [...current.stories, story] }
        : current,
    );
    setCreating(false);
  }
  async function retry(story: UserStory) {
    setRetrying(story.id);
    setError("");
    try {
      const session = await (
        story.sessionKind === "implementation" ? retryImplementation : retryPlanning
      )(story.id);
      setList((current) =>
        current.status === "ready"
          ? {
              status: "ready",
              stories: current.stories.map((item) =>
                item.id === story.id
                  ? {
                      ...item,
                      planningStatus: session.status,
                      sessionId: session.id,
                      sessionKind: session.kind,
                    }
                  : item,
              ),
            }
          : current,
      );
      void navigate(
        `/stories/${encodeURIComponent(story.id)}/sessions/${encodeURIComponent(session.id)}`,
      );
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not start session.");
    } finally {
      setRetrying("");
    }
  }

  return (
    <>
      <header className="page-title">
        <h1>{factory.name}</h1>
      </header>
      <div className="section-toolbar">
        <h2 className="section-title">User stories</h2>
        <Button
          variant="primary"
          shape="square"
          aria-label="New user story"
          title="New user story"
          onClick={() => setCreating(true)}
          disabled={list.status !== "ready" || creating}
        >
          <Plus size={16} />
        </Button>
      </div>
      {error && (
        <p className="error-message" role="alert">
          {error}
        </p>
      )}
      {creating && (
        <UserStoryForm
          factory={factory}
          onCreated={storyCreated}
          onCancel={() => setCreating(false)}
        />
      )}
      {list.status === "loading" ? (
        <output className="empty-state">Loading…</output>
      ) : list.status === "failed" ? (
        <p className="error-message" role="alert">
          Could not load user stories.
        </p>
      ) : (
        <DataTable
          rows={list.stories}
          columns={columns}
          rowKey={(story) => story.id}
          label="User stories"
          emptyMessage="No user stories"
          actions={(story) => (
            <>
              <Button
                variant="ghost"
                size="sm"
                shape="square"
                aria-label={`Open session for ${story.title}`}
                title="Session"
                disabled={!story.sessionId}
                onClick={() => {
                  void navigate(
                    `/stories/${encodeURIComponent(story.id)}/sessions/${encodeURIComponent(story.sessionId)}`,
                  );
                }}
              >
                <Terminal size={16} />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                shape="square"
                aria-label={`Retry planning ${story.title}`}
                title="Plan again"
                disabled={!story.goal || isActive(story.planningStatus) || !!retrying}
                loading={retrying === story.id}
                onClick={() => {
                  void retry(story);
                }}
              >
                <ArrowClockwise size={16} />
              </Button>
            </>
          )}
        />
      )}
    </>
  );
}
