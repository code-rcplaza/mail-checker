import type { ReactNode } from "react";
import { classNames } from "../../../lib/classNames";

export interface ListRowProps {
  readonly title: string;
  readonly subtitle?: string | undefined;
  readonly value?: ReactNode;
  readonly leading?: ReactNode;
  /** With an action the row is a real button: touch, click and Enter all work. */
  readonly onSelect?: (() => void) | undefined;
}

export function ListRow({ title, subtitle, value, leading, onSelect }: ListRowProps) {
  const content = (
    <>
      {leading}
      <span className="flex min-w-0 flex-1 flex-col">
        <span className="truncate font-medium text-ink">{title}</span>
        {subtitle !== undefined && <span className="truncate text-sm text-ink-muted">{subtitle}</span>}
      </span>
      {value !== undefined && <span className="flex-none text-right font-semibold text-ink">{value}</span>}
    </>
  );
  const rowClasses = "flex min-h-11 w-full items-center gap-4 py-2 text-left";
  if (onSelect === undefined) {
    return <div className={rowClasses}>{content}</div>;
  }
  return (
    <button type="button" onClick={onSelect} className={classNames(rowClasses, "rounded-md")}>
      {content}
    </button>
  );
}
