import { createMemoryHistory } from "@tanstack/react-router";
import { render } from "@testing-library/react";
import type { Gateways } from "../lib/gateways/gateways";
import { gatewaysWith } from "../lib/gateways/renderWithGateways";
import { AppProviders } from "./AppProviders";
import { createApp } from "./createApp";

/** The whole app at a path, with the ports a test gives and a memory history; no retries. */
export function renderApp(path: string, given: Partial<Gateways> = {}) {
  const gateways = gatewaysWith(given);
  const { queryClient, router } = createApp(gateways, createMemoryHistory({ initialEntries: [path] }));
  queryClient.setDefaultOptions({ queries: { retry: false }, mutations: { retry: false } });
  const rendered = render(<AppProviders queryClient={queryClient} gateways={gateways} router={router} />);
  return { ...rendered, router, queryClient };
}
