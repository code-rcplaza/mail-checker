import { describe, expect, it } from "vitest";
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
