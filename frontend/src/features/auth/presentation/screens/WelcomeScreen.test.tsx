import { screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { renderApp } from "../../../../app/renderApp";
import { gatewaysWith } from "../../../../lib/gateways/renderWithGateways";

const ana = { id: "u1", name: "Ana Pérez", email: "ana@example.com", initials: "AP" };

describe("Welcome screen", () => {
  it("greets with the name of the session and only reads /api/me", async () => {
    const auth = { ...gatewaysWith({}).auth, getMe: vi.fn(async () => ana) };
    renderApp("/bienvenida", { auth });

    expect(await screen.findByText("Hola, Ana")).toBeInTheDocument();
    expect(auth.getMe).toHaveBeenCalledOnce();
  });

  it("is a screen with a session: avatar and navigation", async () => {
    renderApp("/bienvenida", { auth: { ...gatewaysWith({}).auth, getMe: async () => ana } });

    expect(await screen.findAllByRole("link", { name: "Perfil" })).not.toHaveLength(0);
    expect(screen.getAllByRole("navigation", { name: "Navegación principal" })).not.toHaveLength(0);
    expect(screen.queryByRole("button", { name: /Cambiar a tema/ })).not.toBeInTheDocument();
  });
});
