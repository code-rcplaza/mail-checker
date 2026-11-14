import { createContext, type ReactNode, useContext } from "react";
import type { Gateways } from "./gateways";

const GatewaysContext = createContext<Gateways | null>(null);

export function GatewaysProvider({
  gateways,
  children,
}: {
  readonly gateways: Gateways;
  readonly children: ReactNode;
}) {
  return <GatewaysContext.Provider value={gateways}>{children}</GatewaysContext.Provider>;
}

/** The ports for the hooks of `application`; outside the provider it fails loudly. */
export function useGateways(): Gateways {
  const gateways = useContext(GatewaysContext);
  if (gateways === null) {
    throw new Error("useGateways needs a GatewaysProvider above it");
  }
  return gateways;
}
