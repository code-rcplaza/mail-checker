import { BackIcon } from "../ui/atoms/icons";
import { AppLink } from "./AppLink";
import type { Navigate } from "./navigation";

export interface BackLinkProps {
  readonly to: string;
  /** The name of where it goes back to ("Fuente"), as the handoff shows it. */
  readonly label: string;
  readonly onNavigate: Navigate;
}

export function BackLink({ to, label, onNavigate }: BackLinkProps) {
  return (
    <AppLink
      to={to}
      onNavigate={onNavigate}
      className="inline-flex min-h-11 items-center gap-1 text-sm font-semibold text-zone-ink"
    >
      <BackIcon />
      {label}
    </AppLink>
  );
}
