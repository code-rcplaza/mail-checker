import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { sessionQuery } from "../features/auth/application/sessionQuery";

/** Every screen with the person's data: without a session, back to Login. */
export const Route = createFileRoute("/_session")({
  beforeLoad: async ({ context }) => {
    const me = await context.queryClient.ensureQueryData(sessionQuery(context.gateways.auth));
    if (me === null) {
      throw redirect({ to: "/entrar" });
    }
    return { me };
  },
  component: Outlet,
});
