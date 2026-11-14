import { describe, expect, it, vi } from "vitest";
import { ApiError, NetworkError } from "../api/client";
import { createQueryClient, shouldRetry } from "./queryClient";

describe("query retries", () => {
  it("never retries an API error", () => {
    expect(shouldRetry(0, new ApiError({ status: 404, key: "common.not_found" }))).toBe(false);
    expect(shouldRetry(0, new ApiError({ status: 400, key: "common.validation_failed" }))).toBe(false);
  });

  it("retries a network failure once", () => {
    expect(shouldRetry(0, new NetworkError(new TypeError("Failed to fetch")))).toBe(true);
    expect(shouldRetry(1, new NetworkError(new TypeError("Failed to fetch")))).toBe(false);
  });

  it("is the policy of the app's queries, and mutations never retry", () => {
    const defaults = createQueryClient().getDefaultOptions();

    expect(defaults.queries?.retry).toBe(shouldRetry);
    expect(defaults.mutations?.retry).toBe(false);
  });
});

describe("unauthorized answers", () => {
  const unauthorized = new ApiError({ status: 401, key: "common.unauthorized" });

  it("leave the session when a query gets a 401", async () => {
    const onUnauthorized = vi.fn();
    const client = createQueryClient({ onUnauthorized });

    await client.fetchQuery({ queryKey: ["x"], queryFn: () => Promise.reject(unauthorized) }).catch(() => {});

    expect(onUnauthorized).toHaveBeenCalledOnce();
  });

  it("leave the session when a mutation gets a 401", async () => {
    const onUnauthorized = vi.fn();
    const client = createQueryClient({ onUnauthorized });
    const mutation = client
      .getMutationCache()
      .build(client, { mutationFn: () => Promise.reject(unauthorized) });

    await mutation.execute(undefined).catch(() => {});

    expect(onUnauthorized).toHaveBeenCalledOnce();
  });

  it("ignore any other error", async () => {
    const onUnauthorized = vi.fn();
    const client = createQueryClient({ onUnauthorized });
    const notFound = new ApiError({ status: 404, key: "common.not_found" });

    await client.fetchQuery({ queryKey: ["x"], queryFn: () => Promise.reject(notFound) }).catch(() => {});

    expect(onUnauthorized).not.toHaveBeenCalled();
  });
});
