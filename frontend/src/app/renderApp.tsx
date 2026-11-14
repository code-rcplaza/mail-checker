import { createMemoryHistory } from "@tanstack/react-router";
import { render } from "@testing-library/react";
import type { Gateways } from "../lib/gateways/gateways";
import { gatewaysWith, testQueryClient } from "../lib/gateways/renderWithGateways";
import { AppProviders } from "./AppProviders";
import { createAppRouter } from "./router";

/** The whole app at a path, with the ports a test gives and a memory history. */
export function renderApp(path: string, given: Partial<Gateways> = {}) {
  const queryClient = testQueryClient();
  const gateways = gatewaysWith(given);
  const router = createAppRouter({ queryClient, gateways }, createMemoryHistory({ initialEntries: [path] }));
  const rendered = render(<AppProviders queryClient={queryClient} gateways={gateways} router={router} />);
  return { ...rendered, router };
}
