import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach, vi } from "vitest";

// No test reaches the network: a test that needs the API passes its own fetch double.
const blockedFetch = vi.fn(async (input: RequestInfo | URL) => {
  throw new Error(`fetch is not available in tests: ${String(input)}`);
});

vi.stubGlobal("fetch", blockedFetch);

afterEach(() => {
  cleanup();
  blockedFetch.mockClear();
});
