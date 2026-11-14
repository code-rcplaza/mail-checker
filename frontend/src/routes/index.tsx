import { createFileRoute, redirect } from "@tanstack/react-router";
import { sessionQuery } from "../features/auth/application/sessionQuery";

/**
 * The root is not a screen: with a session it opens the Resumen, without one
 * (or if /api/me cannot be read) Login, which links to Registro.
 */
export const Route = createFileRoute("/")({
  beforeLoad: async ({ context }) => {
    const me = await context.queryClient
      .ensureQueryData(sessionQuery(context.gateways.auth))
      .catch(() => null);
    throw redirect({ to: me === null ? "/entrar" : "/resumen" });
  },
});
