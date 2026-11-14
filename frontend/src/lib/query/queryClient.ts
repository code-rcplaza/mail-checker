import { MutationCache, QueryCache, QueryClient } from "@tanstack/react-query";
import { ApiError, NetworkError } from "../api/client";

/** A 4xx will answer the same again; a network failure may not, so it gets one more try. */
export function shouldRetry(failureCount: number, error: unknown): boolean {
  return error instanceof NetworkError && failureCount < 1;
}

export interface QueryClientOptions {
  /** A 401 from any query or mutation: the session ended on the server. */
  readonly onUnauthorized?: () => void;
}

function isUnauthorized(error: unknown): boolean {
  return error instanceof ApiError && error.isUnauthorized;
}

export function createQueryClient({ onUnauthorized }: QueryClientOptions = {}): QueryClient {
  const leaveOnUnauthorized = (error: unknown) => {
    if (isUnauthorized(error)) {
      onUnauthorized?.();
    }
  };
  return new QueryClient({
    queryCache: new QueryCache({ onError: leaveOnUnauthorized }),
    mutationCache: new MutationCache({ onError: leaveOnUnauthorized }),
    defaultOptions: {
      queries: { retry: shouldRetry, refetchOnWindowFocus: false },
      mutations: { retry: false },
    },
  });
}
