import { QueryClient } from "@tanstack/react-query";
import { NetworkError } from "../api/client";

/** A 4xx will answer the same again; a network failure may not, so it gets one more try. */
export function shouldRetry(failureCount: number, error: unknown): boolean {
  return error instanceof NetworkError && failureCount < 1;
}

export function createQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: shouldRetry, refetchOnWindowFocus: false },
      mutations: { retry: false },
    },
  });
}
