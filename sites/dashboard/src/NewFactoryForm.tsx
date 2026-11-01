import { useCallback, useState, type FormEvent } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { createFactory, type Factory } from "./api";

export function NewFactoryForm({
  onCreated,
  onCancel,
}: {
  onCreated: (factory: Factory) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const focusName = useCallback((element: HTMLInputElement | null) => element?.focus(), []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!name.trim() || saving) return;
    setSaving(true);
    setError("");
    try {
      onCreated(await createFactory(name));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not create factory.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <form
      className="create-form"
      onSubmit={(event) => {
        void submit(event);
      }}
    >
      <label htmlFor="factory-name">Name</label>
      <input
        id="factory-name"
        value={name}
        onChange={(event) => setName(event.target.value)}
        maxLength={80}
        disabled={saving}
        ref={focusName}
      />
      {error && (
        <p className="error-message" role="alert">
          {error}
        </p>
      )}
      <div>
        <Button variant="ghost" size="sm" type="button" onClick={onCancel} disabled={saving}>
          Cancel
        </Button>
        <Button
          variant="primary"
          size="sm"
          type="submit"
          loading={saving}
          disabled={!name.trim() || saving}
        >
          Create
        </Button>
      </div>
    </form>
  );
}
