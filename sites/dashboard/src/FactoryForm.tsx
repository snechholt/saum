import { useCallback, useState, type FormEvent } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { createFactory, updateFactory, type Factory } from "./api";

export function FactoryForm({
  factory,
  onSaved,
  onCancel,
}: {
  factory?: Factory;
  onSaved: (factory: Factory) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState(factory?.name ?? "");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const focusName = useCallback((element: HTMLInputElement | null) => element?.focus(), []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!name.trim() || saving) return;
    setSaving(true);
    setError("");
    try {
      onSaved(await (factory ? updateFactory(factory.id, name) : createFactory(name)));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Could not save factory.");
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
          {factory ? "Save" : "Create"}
        </Button>
      </div>
    </form>
  );
}
