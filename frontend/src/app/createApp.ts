import type { RouterHistory } from "@tanstack/react-router";
import type { Gateways } from "../lib/gateways/gateways";
import { createQueryClient } from "../lib/query/queryClient";
import { createAppRouter } from "./router";

/**
 * The query client and the router, wired: a 401 from any call drops every
 * cached answer and goes back to Login. main.tsx and the tests build the same.
 */
export function createApp(gateways: Gateways, history?: RouterHistory) {
  const leaveSession = () => {
    queryClient.clear();
    void router.navigate({ to: "/entrar" });
  };
  const queryClient = createQueryClient({ onUnauthorized: leaveSession });
  const router = createAppRouter({ queryClient, gateways }, history);
  return { queryClient, router };
}
