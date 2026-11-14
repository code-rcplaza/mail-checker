import { queryOptions } from "@tanstack/react-query";
import type { AuthGateway } from "./ports";

export const sessionQueryKey = ["session"] as const;

/** `GET /api/me` in the cache: the guards and `useSession` read the same entry. */
export function sessionQuery(auth: AuthGateway) {
  return queryOptions({ queryKey: sessionQueryKey, queryFn: () => auth.getMe(), staleTime: 5 * 60 * 1000 });
}
