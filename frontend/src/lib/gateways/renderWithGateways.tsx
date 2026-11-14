import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type RenderResult, render } from "@testing-library/react";
import type { ReactElement } from "react";
import type { AuthGateway } from "../../features/auth/application/ports";
import type { StatementsGateway } from "../../features/statements/application/ports";
import { GatewaysProvider } from "./GatewaysProvider";
import type { Gateways } from "./gateways";

function notGiven(name: string): Promise<never> {
  return Promise.reject(new Error(`the ${name} gateway was not given to this test`));
}

const missingAuth: AuthGateway = {
  getMe: () => notGiven("auth"),
  register: () => notGiven("auth"),
  login: () => notGiven("auth"),
  verifyCode: () => notGiven("auth"),
  resendCode: () => notGiven("auth"),
  logout: () => notGiven("auth"),
};

/** Fills the ports a test did not give with doubles that fail when used. */
const missingStatements: StatementsGateway = {
  listBanks: () => notGiven("statements"),
  upload: () => notGiven("statements"),
};

export function gatewaysWith(given: Partial<Gateways>): Gateways {
  return { auth: given.auth ?? missingAuth, statements: given.statements ?? missingStatements };
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
