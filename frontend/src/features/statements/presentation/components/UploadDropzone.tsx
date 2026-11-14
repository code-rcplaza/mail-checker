import { type ChangeEvent, type DragEvent, useRef } from "react";
import { Button } from "../../../../components/ui/atoms/Button";
import { t } from "../../../../lib/i18n/translate";

export interface UploadDropzoneProps {
  readonly file: File | undefined;
  /** Already translated: why the last file was rejected. */
  readonly error: string | undefined;
  readonly onPick: (file: File) => void;
}

/**
 * Drop a PDF or choose one. The real control is a button (Enter and space
 * work by themselves) that opens the hidden file input.
 */
export function UploadDropzone({ file, error, onPick }: UploadDropzoneProps) {
  const input = useRef<HTMLInputElement>(null);

  function pickFromInput(event: ChangeEvent<HTMLInputElement>) {
    const picked = event.target.files?.[0];
    // Cleared so choosing the same file again still counts as a change.
    event.target.value = "";
    if (picked !== undefined) {
      onPick(picked);
    }
  }

  function pickFromDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    const dropped = event.dataTransfer.files[0];
    if (dropped !== undefined) {
      onPick(dropped);
    }
  }

  return (
    <div className="flex flex-col gap-2">
      {/* biome-ignore lint/a11y/noStaticElementInteractions: dropping is a mouse shortcut; the button inside is the keyboard way */}
      <div
        data-testid="dropzone"
        onDragOver={(event) => event.preventDefault()}
        onDrop={pickFromDrop}
        className="flex flex-col items-center gap-3 rounded-lg border-[1.5px] border-dashed border-border-strong bg-field px-4 py-6 text-center"
      >
        <input
          ref={input}
          type="file"
          accept="application/pdf,.pdf"
          aria-label={t("ui.upload.file")}
          tabIndex={-1}
          onChange={pickFromInput}
          className="sr-only"
        />
        <Button variant="secondary" onClick={() => input.current?.click()}>
          {file === undefined ? t("ui.upload.pick") : t("ui.upload.change")}
        </Button>
        {file === undefined ? (
          <p className="text-sm text-ink-muted">{t("ui.upload.drop")}</p>
        ) : (
          <p className="max-w-full truncate text-sm font-semibold text-ink">{file.name}</p>
        )}
      </div>
      {error !== undefined && (
        <p role="alert" className="text-sm text-danger-ink">
          {error}
        </p>
      )}
    </div>
  );
}
