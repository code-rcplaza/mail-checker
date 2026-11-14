import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { NotFoundScreen } from "../app/screens/NotFoundScreen";
import type { Gateways } from "../lib/gateways/gateways";

/** What every route can reach without importing infrastructure. */
export interface RouterContext {
  readonly queryClient: QueryClient;
  readonly gateways: Gateways;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: Outlet,
  notFoundComponent: NotFoundScreen,
});
