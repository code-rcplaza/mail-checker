import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { sessionQuery } from "../features/auth/application/sessionQuery";

/**
 * Registro and Login are for people without a session. If /api/me cannot be
 * read (offline), they still open: entering is how one gets a session back.
 */
export const Route = createFileRoute("/_guest")({
  beforeLoad: async ({ context }) => {
    const me = await context.queryClient
      .ensureQueryData(sessionQuery(context.gateways.auth))
      .catch(() => null);
    if (me !== null) {
      throw redirect({ to: "/resumen" });
    }
  },
  component: Outlet,
});
