import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type RenderResult, render } from "@testing-library/react";
import type { ReactElement } from "react";
import type { AuthGateway } from "../../features/auth/application/ports";
import { GatewaysProvider } from "./GatewaysProvider";
import type { Gateways } from "./gateways";

function notGiven(name: string): Promise<never> {
  return Promise.reject(new Error(`the ${name} gateway was not given to this test`));
}

const missingAuth: AuthGateway = { getMe: () => notGiven("auth") };

/** Fills the ports a test did not give with doubles that fail when used. */
export function gatewaysWith(given: Partial<Gateways>): Gateways {
  return { auth: given.auth ?? missingAuth };
}

/** A fresh QueryClient per test, without retries, so tests never share cache nor wait. */
export function testQueryClient(): QueryClient {
  return new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
}

export function renderWithGateways(
  ui: ReactElement,
  given: Partial<Gateways> = {},
): RenderResult & { queryClient: QueryClient } {
  const queryClient = testQueryClient();
  const rendered = render(
    <QueryClientProvider client={queryClient}>
      <GatewaysProvider gateways={gatewaysWith(given)}>{ui}</GatewaysProvider>
    </QueryClientProvider>,
  );
  return { ...rendered, queryClient };
}
