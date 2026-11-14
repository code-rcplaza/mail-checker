import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { NetworkError } from "../../../lib/api/client";
import { gatewaysWith, renderWithGateways } from "../../../lib/gateways/renderWithGateways";
import type { AuthGateway } from "./ports";
import { useLogout } from "./useLogout";
import { useSession } from "./useSession";

const ana = { id: "u1", name: "Ana Pérez", email: "ana@example.com", initials: "AP" };

function authWith(overrides: Partial<AuthGateway>): AuthGateway {
  return { ...gatewaysWith({}).auth, ...overrides };
}

function SessionState() {
  const session = useSession();
  if (session.status === "signedIn") {
    return <p>{`${session.me.name} · ${session.me.email} · ${session.me.initials}`}</p>;
  }
  return <p>{session.status}</p>;
}

describe("useSession", () => {
  it("brings name, email and initials with a session", async () => {
    renderWithGateways(<SessionState />, { auth: authWith({ getMe: async () => ana }) });

    expect(await screen.findByText("Ana Pérez · ana@example.com · AP")).toBeInTheDocument();
  });

  it("is no session, not an error, without one", async () => {
    renderWithGateways(<SessionState />, { auth: authWith({ getMe: async () => null }) });

    expect(await screen.findByText("signedOut")).toBeInTheDocument();
  });

  it("is loading until /api/me answers, never signed out first", () => {
    renderWithGateways(<SessionState />, { auth: authWith({ getMe: () => new Promise(() => {}) }) });

    expect(screen.getByText("loading")).toBeInTheDocument();
  });
});

function LogoutButton({ onLoggedOut }: { readonly onLoggedOut: (reachedServer: boolean) => void }) {
  const { logout } = useLogout({ onLoggedOut });
  return (
    <button type="button" onClick={logout}>
      salir
    </button>
  );
}

describe("useLogout", () => {
  it("calls logout, empties the cache and lets the screen go to Login", async () => {
    const auth = authWith({ logout: vi.fn(async () => {}) });
    const onLoggedOut = vi.fn();
    const { queryClient } = renderWithGateways(<LogoutButton onLoggedOut={onLoggedOut} />, { auth });
    queryClient.setQueryData(["session"], ana);

    await userEvent.click(screen.getByRole("button", { name: "salir" }));

    await waitFor(() => expect(onLoggedOut).toHaveBeenCalledWith(true));
    expect(auth.logout).toHaveBeenCalledOnce();
    expect(queryClient.getQueryCache().getAll()).toHaveLength(0);
  });

  it("still cleans what is local when the network fails, and says the server was not reached", async () => {
    const auth = authWith({
      logout: vi.fn(async () => {
        throw new NetworkError(new TypeError("Failed to fetch"));
      }),
    });
    const onLoggedOut = vi.fn();
    const { queryClient } = renderWithGateways(<LogoutButton onLoggedOut={onLoggedOut} />, { auth });
    queryClient.setQueryData(["session"], ana);

    await userEvent.click(screen.getByRole("button", { name: "salir" }));

    await waitFor(() => expect(onLoggedOut).toHaveBeenCalledWith(false));
    expect(queryClient.getQueryCache().getAll()).toHaveLength(0);
  });
});
