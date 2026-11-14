import { createFileRoute } from "@tanstack/react-router";

/** Placeholder until the Resumen screen (FOKKO-0103-SUBS). */
function SummaryRoute() {
  return <main />;
}

export const Route = createFileRoute("/_session/resumen")({ component: SummaryRoute });
