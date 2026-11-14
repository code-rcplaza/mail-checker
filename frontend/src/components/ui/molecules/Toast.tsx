import { useEffect } from "react";
import { t } from "../../../lib/i18n/translate";

export const toastDuration = 4000;

export interface ToastProps {
  readonly message: string;
  readonly onClose: () => void;
}

/**
 * One notice at a time: the screen keeps a single toast in state, so a new
 * one replaces the old (and, rendered with a new `key`, restarts the timer).
 */
export function Toast({ message, onClose }: ToastProps) {
  useEffect(() => {
    const timer = window.setTimeout(onClose, toastDuration);
    const closeOnEscape = (event: globalThis.KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };
    document.addEventListener("keydown", closeOnEscape);
    return () => {
      window.clearTimeout(timer);
      document.removeEventListener("keydown", closeOnEscape);
    };
  }, [onClose]);

  return (
    <div
      role="status"
      aria-live="polite"
      className="fixed inset-x-4 bottom-24 z-50 mx-auto flex max-w-md items-center gap-4 rounded-lg bg-ink px-4 py-3 text-page lg:bottom-8"
    >
      <p className="flex-1 text-sm">{message}</p>
      <button type="button" onClick={onClose} className="min-h-11 rounded-md px-2 text-sm font-semibold">
        {t("ui.close")}
      </button>
    </div>
  );
}
