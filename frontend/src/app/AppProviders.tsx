import { type QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider } from "@tanstack/react-router";
import { GatewaysProvider } from "../lib/gateways/GatewaysProvider";
import type { Gateways } from "../lib/gateways/gateways";
import type { createAppRouter } from "./router";

export interface AppProvidersProps {
  readonly queryClient: QueryClient;
  readonly gateways: Gateways;
  readonly router: ReturnType<typeof createAppRouter>;
}

/** Server data lives only in TanStack Query; the ports come by Context. */
export function AppProviders({ queryClient, gateways, router }: AppProvidersProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <GatewaysProvider gateways={gateways}>
        <RouterProvider router={router} />
      </GatewaysProvider>
    </QueryClientProvider>
  );
}
