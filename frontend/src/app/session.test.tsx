import { screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { AuthGateway } from "../features/auth/application/ports";
import { ApiError } from "../lib/api/client";
import { gatewaysWith } from "../lib/gateways/renderWithGateways";
import { renderApp } from "./renderApp";

const ana = { id: "u1", name: "Ana Pérez", email: "ana@example.com", initials: "AP" };

function authReading(me: typeof ana | null): AuthGateway {
  return { ...gatewaysWith({}).auth, getMe: vi.fn(async () => me) };
}

describe("route guards", () => {
  it("sends a protected route without a session to Login", async () => {
    const { router } = renderApp("/resumen", { auth: authReading(null) });

    await waitFor(() => expect(router.state.location.pathname).toBe("/entrar"));
  });

  it.each(["/entrar", "/registro"])("sends %s with a session to the Resumen", async (path) => {
    const { router } = renderApp(path, { auth: authReading(ana) });

    await waitFor(() => expect(router.state.location.pathname).toBe("/resumen"));
  });

  it("opens a protected route with a session", async () => {
    const { router } = renderApp("/resumen", { auth: authReading(ana) });

    await waitFor(() => expect(router.state.status).toBe("idle"));
    expect(router.state.location.pathname).toBe("/resumen");
  });

  it("goes back to Login and drops the cache on a 401 from any call", async () => {
    // The server ended the session: after the 401, /api/me answers without one too.
    const auth = { ...authReading(ana), getMe: vi.fn().mockResolvedValueOnce(ana).mockResolvedValue(null) };
    const { router, queryClient } = renderApp("/resumen", { auth });
    await waitFor(() => expect(router.state.location.pathname).toBe("/resumen"));
    queryClient.setQueryData(["private"], "datos");

    await queryClient
      .fetchQuery({
        queryKey: ["expired"],
        queryFn: () => Promise.reject(new ApiError({ status: 401, key: "common.unauthorized" })),
      })
      .catch(() => {});

    await waitFor(() => expect(router.state.location.pathname).toBe("/entrar"));
    expect(queryClient.getQueryData(["private"])).toBeUndefined();
    expect(await screen.findByRole("heading", { name: "Entra a Fokko" })).toBeInTheDocument();
  });
});
