import { useQuery } from "@tanstack/react-query";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import type { Me } from "../domain/me";
import { sessionQuery } from "./sessionQuery";

/** Loading is its own state, so Login never flashes while /api/me answers. */
export type Session =
  | { readonly status: "loading" }
  | { readonly status: "signedOut" }
  | { readonly status: "signedIn"; readonly me: Me };

export function useSession(): Session {
  const { auth } = useGateways();
  const reading = useQuery(sessionQuery(auth));
  if (reading.isPending) {
    return { status: "loading" };
  }
  if (reading.data === null || reading.data === undefined) {
    return { status: "signedOut" };
  }
  return { status: "signedIn", me: reading.data };
}
