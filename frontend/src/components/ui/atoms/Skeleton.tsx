import { classNames } from "../../../lib/classNames";

/** A placeholder while data loads; screen readers skip it, the loading state is announced elsewhere. */
export function Skeleton({ className }: { readonly className?: string }) {
  return (
    <span
      aria-hidden="true"
      className={classNames("block rounded-md bg-skeleton motion-safe:animate-pulse", className)}
    />
  );
}
