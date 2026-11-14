import { describe, expect, it, vi } from "vitest";
import { ApiError, createApiClient, NetworkError } from "./client";

function respondWith(body: BodyInit | null, init: ResponseInit) {
  return vi.fn(async (_input: RequestInfo | URL, _init?: RequestInit) => new Response(body, init));
}

describe("api client", () => {
  it("returns the JSON of a 2xx response", async () => {
    const fetchDouble = respondWith(JSON.stringify({ status: "ok" }), { status: 200 });

    await expect(createApiClient(fetchDouble).get("/api/health")).resolves.toEqual({ status: "ok" });
  });

  it("returns nothing for a 204 without reading a body", async () => {
    const fetchDouble = respondWith(null, { status: 204 });

    await expect(
      createApiClient(fetchDouble).send({ method: "POST", path: "/api/auth/logout" }),
    ).resolves.toBeUndefined();
  });

  it("throws a typed error with key, params, fields and status", async () => {
    const body = {
      error: {
        key: "common.validation_failed",
        params: { limit: 3 },
        fields: { email: "auth.email_invalid" },
      },
    };
    const fetchDouble = respondWith(JSON.stringify(body), { status: 400 });

    const error = await createApiClient(fetchDouble)
      .send({ method: "POST", path: "/api/auth/login", body: {} })
      .catch((caught: unknown) => caught);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      status: 400,
      key: "common.validation_failed",
      params: { limit: 3 },
      fields: { email: "auth.email_invalid" },
      isUnauthorized: false,
    });
  });

  it("marks a 401 as unauthorized", async () => {
    const fetchDouble = respondWith(JSON.stringify({ error: { key: "common.unauthorized" } }), {
      status: 401,
    });

    await expect(createApiClient(fetchDouble).get("/api/me")).rejects.toMatchObject({
      isUnauthorized: true,
      params: {},
      fields: {},
    });
  });

  it.each([
    ["not JSON", "<html>bad gateway</html>", 502],
    ["truncated JSON", '{"status":', 200],
    ["JSON without the error shape", '{"message":"oops"}', 500],
  ])("turns %s into common.internal", async (_name, body, status) => {
    const fetchDouble = respondWith(body, { status });

    await expect(createApiClient(fetchDouble).get("/api/summary")).rejects.toMatchObject({
      key: "common.internal",
      status,
    });
  });

  it("tells a network failure apart, as retryable", async () => {
    const fetchDouble = vi.fn(async () => {
      throw new TypeError("Failed to fetch");
    });

    const error = await createApiClient(fetchDouble)
      .get("/api/summary")
      .catch((caught: unknown) => caught);

    expect(error).toBeInstanceOf(NetworkError);
    expect(error).not.toBeInstanceOf(ApiError);
    expect(error).toMatchObject({ isRetryable: true });
  });

  it("sends the session cookie, same origin, and JSON for writes", async () => {
    const fetchDouble = respondWith(null, { status: 204 });

    await createApiClient(fetchDouble).send({
      method: "PUT",
      path: "/api/subscriptions/aliases",
      body: { lineKey: "k", alias: "Casa" },
    });

    const [path, init] = fetchDouble.mock.calls[0] ?? [];
    expect(path).toBe("/api/subscriptions/aliases");
    expect(init).toMatchObject({
      method: "PUT",
      credentials: "same-origin",
      mode: "same-origin",
      body: '{"lineKey":"k","alias":"Casa"}',
    });
    expect(new Headers(init?.headers).get("Content-Type")).toBe("application/json");
  });

  it("sends a FormData as multipart, letting the browser set the boundary", async () => {
    const fetchDouble = respondWith(JSON.stringify({ id: "s1" }), { status: 201 });
    const form = new FormData();
    form.append("bank", "bci");

    await createApiClient(fetchDouble).send({ method: "POST", path: "/api/statements", body: form });

    const init = fetchDouble.mock.calls[0]?.[1];
    expect(init?.body).toBe(form);
    expect(init?.headers).toBeUndefined();
  });
});
