import { type MessageKey, t } from "../../../lib/i18n/translate";

export type BadgeKind = "toCancel" | "rise" | "possible";

const badges: Readonly<Record<BadgeKind, { label: MessageKey; classes: string }>> = {
  toCancel: { label: "ui.badge.to_cancel", classes: "bg-error-surface text-danger-ink" },
  rise: { label: "ui.badge.rise", classes: "bg-attention-solid text-on-attention-solid" },
  possible: { label: "ui.badge.possible", classes: "bg-nav-active text-primary" },
};

/** A state of a line; it always says it in words, never only with color. */
export function Badge({ kind }: { readonly kind: BadgeKind }) {
  const badge = badges[kind];
  return (
    <span
      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-semibold ${badge.classes}`}
    >
      {t(badge.label)}
    </span>
  );
}
