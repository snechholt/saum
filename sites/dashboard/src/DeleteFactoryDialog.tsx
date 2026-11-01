import { Button } from "@cloudflare/kumo/components/button";
import { Dialog } from "@cloudflare/kumo/components/dialog";
import { Trash, X } from "@phosphor-icons/react";
import type { Factory } from "./api";

export function DeleteFactoryDialog({
  factory,
  busy,
  onCancel,
  onConfirm,
}: {
  factory: Factory | null;
  busy: boolean;
  onCancel: () => void;
  onConfirm: (factory: Factory) => void;
}) {
  return (
    <Dialog.Root
      role="alertdialog"
      open={factory !== null}
      onOpenChange={(open) => {
        if (!open && !busy) onCancel();
      }}
    >
      <Dialog size="sm" className="confirm-dialog">
        <Dialog.Title>Delete factory?</Dialog.Title>
        <Dialog.Description>{factory?.name}</Dialog.Description>
        <div className="form-actions">
          <Button
            variant="secondary"
            shape="square"
            aria-label="Cancel deletion"
            title="Cancel"
            disabled={busy}
            onClick={onCancel}
          >
            <X size={16} />
          </Button>
          <Button
            variant="destructive"
            shape="square"
            aria-label="Confirm deletion"
            title="Delete"
            disabled={busy || !factory}
            loading={busy}
            onClick={() => {
              if (factory) onConfirm(factory);
            }}
          >
            <Trash size={16} />
          </Button>
        </div>
      </Dialog>
    </Dialog.Root>
  );
}
