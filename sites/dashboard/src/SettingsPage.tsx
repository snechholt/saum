import { useState } from "react";
import { Button } from "@cloudflare/kumo/components/button";
import { Input } from "@cloudflare/kumo/components/input";
import { LayerCard } from "@cloudflare/kumo/components/layer-card";
import { Select } from "@cloudflare/kumo/components/select";
import { Check, FloppyDisk } from "@phosphor-icons/react";
import { useSettings } from "./settings";
import type { Settings, Theme } from "./settings-api";
function SettingsForm({
  initial,
  save,
}: {
  initial: Settings;
  save: (settings: Settings) => Promise<void>;
}) {
  const [draft, setDraft] = useState(initial),
    [busy, setBusy] = useState(false),
    [error, setError] = useState(""),
    [saved, setSaved] = useState(false);
  return (
    <LayerCard className="settings-surface">
      <form
        className="create-form"
        onSubmit={(event) => {
          event.preventDefault();
          setBusy(true);
          setError("");
          void save(draft)
            .then(() => setSaved(true))
            .catch((reason) =>
              setError(reason instanceof Error ? reason.message : "Could not save settings."),
            )
            .finally(() => setBusy(false));
        }}
      >
        <div className="form-fields">
          <Input
            label="Factories directory"
            value={draft.factoriesRoot}
            required
            disabled={busy}
            onChange={(event) => {
              setSaved(false);
              setDraft({ ...draft, factoriesRoot: event.target.value });
            }}
          />
          <p className="settings-note">
            Applies to new factories. Existing workspaces stay where they are.
          </p>
          <Select<Theme>
            label="Theme"
            value={draft.theme}
            disabled={busy}
            items={[
              { value: "system", label: "System" },
              { value: "light", label: "Light" },
              { value: "dark", label: "Dark" },
            ]}
            onValueChange={(theme) => {
              if (theme) {
                setSaved(false);
                setDraft({ ...draft, theme });
              }
            }}
          />
        </div>
        {error && (
          <p className="error-message" role="alert">
            {error}
          </p>
        )}
        <div className="form-actions">
          <Button
            variant="primary"
            shape="square"
            type="submit"
            title={saved ? "Saved" : "Save settings"}
            aria-label="Save settings"
            disabled={busy || !draft.factoriesRoot.trim()}
            loading={busy}
          >
            {saved ? <Check size={16} /> : <FloppyDisk size={16} />}
          </Button>
        </div>
      </form>
    </LayerCard>
  );
}
export default function SettingsPage() {
  const { settings, error, save } = useSettings();
  return (
    <>
      <header className="page-title">
        <h1>Settings</h1>
      </header>
      {error ? (
        <p className="error-message" role="alert">
          {error}
        </p>
      ) : settings ? (
        <SettingsForm initial={settings} save={save} />
      ) : (
        <output>Loading…</output>
      )}
    </>
  );
}
