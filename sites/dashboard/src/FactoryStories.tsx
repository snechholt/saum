import { useEffect, useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Plus } from "@phosphor-icons/react";
import { DataTable, type Column } from "@saum/ui";
import { listStories, type Factory, type UserStory } from "./api";
import { UserStoryForm } from "./UserStoryForm";

const columns: Column<UserStory>[] = [
  {
    id: "title",
    header: "Title",
    value: (story) => story.title,
    render: (story) => <span className="story-title">{story.title}</span>,
  },
  {
    id: "goal",
    header: "Goal",
    value: (story) => story.goal,
    render: (story) => <span className="story-goal">{story.goal || "—"}</span>,
  },
];
type StoryList =
  | { status: "loading" }
  | { status: "ready"; stories: UserStory[] }
  | { status: "failed" };

export function FactoryStories({ factory }: { factory: Factory }) {
  const [list, setList] = useState<StoryList>({ status: "loading" });
  const [creating, setCreating] = useState(false);
  useEffect(() => {
    let ignore = false;
    listStories(factory.id)
      .then((stories) => {
        if (!ignore) setList({ status: "ready", stories });
      })
      .catch(() => {
        if (!ignore) setList({ status: "failed" });
      });
    return () => {
      ignore = true;
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
        />
      )}
    </>
  );
}
