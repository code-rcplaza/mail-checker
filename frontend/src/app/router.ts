import { createRouter, type RouterHistory } from "@tanstack/react-router";
import type { RouterContext } from "../routes/__root";
import { routeTree } from "../routeTree.gen";
import { NotFoundScreen } from "./screens/NotFoundScreen";

export function createAppRouter(context: RouterContext, history?: RouterHistory) {
  return createRouter({
    routeTree,
    context,
    defaultNotFoundComponent: NotFoundScreen,
    ...(history === undefined ? {} : { history }),
  });
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof createAppRouter>;
  }
}
