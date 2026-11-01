import { useCallback, useState, type FormEvent } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Input, InputArea } from "@cloudflare/kumo/components/input";
import { LayerCard } from "@cloudflare/kumo/components/layer-card";
import { Plus, X } from "@phosphor-icons/react";
import { createStory, type Factory, type UserStory } from "./api";

export function UserStoryForm({
  factory,
  onCreated,
  onCancel,
}: {
  factory: Factory;
  onCreated: (story: UserStory) => void;
  onCancel: () => void;
}) {
  const [title, setTitle] = useState("");
  const [goal, setGoal] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const focusTitle = useCallback((element: HTMLInputElement | null) => element?.focus(), []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (saving || !title.trim() || !goal.trim()) return;
    setSaving(true);
    setError("");
    try {
      onCreated(await createStory(factory.id, { title, goal }));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not create user story.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <LayerCard className="form-surface story-form-surface">
      <form
        className="create-form"
        onSubmit={(event) => {
          void submit(event);
        }}
      >
        <div className="form-fields">
          <Input
            label="Title"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            maxLength={80}
            disabled={saving}
            ref={focusTitle}
            required
          />
          <InputArea
            label="Goal"
            value={goal}
            onChange={(event) => setGoal(event.target.value)}
            placeholder="What should the user be able to achieve?"
            rows={4}
            maxLength={4000}
            disabled={saving}
            required
          />
        </div>
        {error && (
          <p role="alert" className="error-message">
            {error}
          </p>
        )}
        <div className="form-actions">
          <Button
            variant="secondary"
            shape="square"
            type="button"
            aria-label="Cancel"
            title="Cancel"
            disabled={saving}
            onClick={onCancel}
          >
            <X size={16} />
          </Button>
          <Button
            variant="primary"
            shape="square"
            type="submit"
            aria-label="Create user story"
            title="Create user story"
            loading={saving}
            disabled={saving || !title.trim() || !goal.trim()}
          >
            <Plus size={16} />
          </Button>
        </div>
      </form>
    </LayerCard>
  );
}
