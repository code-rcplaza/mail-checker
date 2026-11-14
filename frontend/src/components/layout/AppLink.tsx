import type { AnchorHTMLAttributes, MouseEvent } from "react";
import type { Navigate } from "./navigation";

export interface AppLinkProps extends Omit<AnchorHTMLAttributes<HTMLAnchorElement>, "href" | "onClick"> {
  readonly to: string;
  readonly onNavigate: Navigate;
}

/** A click with a modifier or another button opens a new tab: that one is the browser's. */
export function isPlainClick(event: MouseEvent<HTMLAnchorElement>): boolean {
  return event.button === 0 && !event.metaKey && !event.ctrlKey && !event.shiftKey && !event.altKey;
}

/** A real link (it has its href) whose plain click is handed to the screen. */
export function AppLink({ to, onNavigate, children, ...props }: AppLinkProps) {
  function followInApp(event: MouseEvent<HTMLAnchorElement>) {
    if (!isPlainClick(event)) {
      return;
    }
    event.preventDefault();
    onNavigate(to);
  }

  return (
    <a {...props} href={to} onClick={followInApp}>
      {children}
    </a>
  );
}
