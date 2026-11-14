import { describe, expect, it, vi } from "vitest";
import { type ApiClient, ApiError, createApiClient } from "../../../lib/api/client";
import { createHttpAuthGateway } from "./httpAuthGateway";

function apiAnswering(body: unknown): ApiClient {
  return { get: async () => body, send: async () => undefined };
}

function fetchAnswering(status: number, body?: unknown) {
  return vi.fn(async () => new Response(body === undefined ? null : JSON.stringify(body), { status }));
}

const registration = { name: "Ana Prueba", email: "ana@example.com", acceptedTerms: true };

describe("http auth gateway", () => {
  it("reads me from the API", async () => {
    const me = { id: "u1", name: "Ana Prueba", email: "ana@example.com", initials: "AP" };

    await expect(createHttpAuthGateway(apiAnswering(me)).getMe()).resolves.toEqual(me);
  });

  it("rejects an answer with another shape", async () => {
    await expect(createHttpAuthGateway(apiAnswering({ id: 1 })).getMe()).rejects.toThrow();
  });

  it("registers with POST /api/auth/register and resolves empty on 204", async () => {
    const fetchFunction = fetchAnswering(204);

    await expect(
      createHttpAuthGateway(createApiClient(fetchFunction)).register(registration),
    ).resolves.toBeUndefined();
    expect(fetchFunction).toHaveBeenCalledWith(
      "/api/auth/register",
      expect.objectContaining({ method: "POST", body: JSON.stringify(registration) }),
    );
  });

  it.each([
    [400, { error: { key: "common.validation_failed", fields: { email: "auth.email_invalid" } } }],
    [409, { error: { key: "common.conflict" } }],
  ])("throws ApiError with the key and the fields on %d", async (status, body) => {
    const gateway = createHttpAuthGateway(createApiClient(fetchAnswering(status, body)));

    const error = await gateway.register(registration).catch((caught: unknown) => caught);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      status,
      key: body.error.key,
      fields: "fields" in body.error ? body.error.fields : {},
    });
  });

  it("asks for a login code with POST /api/auth/login and resolves empty on 204", async () => {
    const fetchFunction = fetchAnswering(204);

    await expect(
      createHttpAuthGateway(createApiClient(fetchFunction)).login("ana@example.com"),
    ).resolves.toBeUndefined();
    expect(fetchFunction).toHaveBeenCalledWith(
      "/api/auth/login",
      expect.objectContaining({ method: "POST", body: JSON.stringify({ email: "ana@example.com" }) }),
    );
  });

  it("throws ApiError with the key when login gets a 400", async () => {
    const body = { error: { key: "common.validation_failed", fields: { email: "auth.email_invalid" } } };
    const gateway = createHttpAuthGateway(createApiClient(fetchAnswering(400, body)));

    await expect(gateway.login("ana@")).rejects.toMatchObject({
      status: 400,
      key: "common.validation_failed",
    });
  });

  it.each([
    [
      "verifies the code",
      "/api/auth/code/verify",
      (gateway: ReturnType<typeof createHttpAuthGateway>) =>
        gateway.verifyCode({ email: "ana@example.com", code: "123456" }),
      { email: "ana@example.com", code: "123456" },
    ],
    [
      "resends the code",
      "/api/auth/code/resend",
      (gateway: ReturnType<typeof createHttpAuthGateway>) => gateway.resendCode("ana@example.com"),
      { email: "ana@example.com" },
    ],
  ])("%s with POST %s and resolves empty on 204", async (_case, path, call, body) => {
    const fetchFunction = fetchAnswering(204);

    await expect(call(createHttpAuthGateway(createApiClient(fetchFunction)))).resolves.toBeUndefined();
    expect(fetchFunction).toHaveBeenCalledWith(
      path,
      expect.objectContaining({ method: "POST", body: JSON.stringify(body) }),
    );
  });

  it("throws ApiError when verifying gets a 400", async () => {
    const body = { error: { key: "auth.code_invalid" } };
    const gateway = createHttpAuthGateway(createApiClient(fetchAnswering(400, body)));

    await expect(gateway.verifyCode({ email: "ana@example.com", code: "000000" })).rejects.toMatchObject({
      status: 400,
      key: "auth.code_invalid",
    });
  });

  it("reads no session from a 401 of /api/me", async () => {
    const body = { error: { key: "common.unauthorized" } };

    await expect(
      createHttpAuthGateway(createApiClient(fetchAnswering(401, body))).getMe(),
    ).resolves.toBeNull();
  });

  it("still throws any other error of /api/me", async () => {
    const gateway = createHttpAuthGateway(
      createApiClient(fetchAnswering(500, { error: { key: "common.internal" } })),
    );

    await expect(gateway.getMe()).rejects.toBeInstanceOf(ApiError);
  });

  it("logs out with POST /api/auth/logout and resolves empty on 204", async () => {
    const fetchFunction = fetchAnswering(204);

    await expect(createHttpAuthGateway(createApiClient(fetchFunction)).logout()).resolves.toBeUndefined();
    expect(fetchFunction).toHaveBeenCalledWith(
      "/api/auth/logout",
      expect.objectContaining({ method: "POST" }),
    );
  });
});
