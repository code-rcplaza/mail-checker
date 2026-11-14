import { describe, expect, it } from "vitest";
import type { ApiClient } from "../../../lib/api/client";
import { createHttpAuthGateway } from "./httpAuthGateway";

function apiAnswering(body: unknown): ApiClient {
  return { get: async () => body, send: async () => undefined };
}

describe("http auth gateway", () => {
  it("reads me from the API", async () => {
    const me = { id: "u1", name: "Ana Prueba", email: "ana@example.com", initials: "AP" };

    await expect(createHttpAuthGateway(apiAnswering(me)).getMe()).resolves.toEqual(me);
  });

  it("rejects an answer with another shape", async () => {
    await expect(createHttpAuthGateway(apiAnswering({ id: 1 })).getMe()).rejects.toThrow();
  });
});
