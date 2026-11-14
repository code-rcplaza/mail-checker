import type { ReactNode } from "react";

export interface EmptyStateProps {
  readonly title: string;
  readonly text: string;
  readonly action?: ReactNode;
}

export function EmptyState({ title, text, action }: EmptyStateProps) {
  return (
    <section className="flex flex-col items-center gap-4 px-6 py-8 text-center">
      <h2 className="text-lg font-semibold text-ink">{title}</h2>
      <p className="text-ink-secondary">{text}</p>
      {action}
    </section>
  );
}
