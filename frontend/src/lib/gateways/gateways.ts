import type { AuthGateway } from "../../features/auth/application/ports";

/** Every port of the front; `main.tsx` builds the real ones, tests give doubles. */
export interface Gateways {
  readonly auth: AuthGateway;
}
