import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { renderApp } from "../../../../app/renderApp";
import { ApiError } from "../../../../lib/api/client";
import { gatewaysWith } from "../../../../lib/gateways/renderWithGateways";
import type { AuthGateway } from "../../application/ports";

function authLoggingIn(login: AuthGateway["login"]): AuthGateway {
  return { ...gatewaysWith({}).auth, login: vi.fn(login) };
}

async function sendEmail(email: string) {
  await userEvent.type(await screen.findByLabelText("Tu correo"), email);
  await userEvent.click(screen.getByRole("button", { name: "Continuar" }));
}

describe("Login screen", () => {
  it("asks for a code and goes to Código with the email", async () => {
    const auth = authLoggingIn(async () => {});
    const { router } = renderApp("/entrar", { auth });

    await sendEmail(" ana@example.com ");

    await waitFor(() => expect(router.state.location.pathname).toBe("/codigo"));
    expect(router.state.location.search).toEqual({ email: "ana@example.com", flow: "login" });
    expect(auth.login).toHaveBeenCalledWith("ana@example.com");
  });

  it("checks the email before calling the server", async () => {
    const auth = authLoggingIn(async () => {});
    renderApp("/entrar", { auth });

    await sendEmail("ana@");

    expect(await screen.findByText("Escribe un correo válido.")).toBeInTheDocument();
    expect(auth.login).not.toHaveBeenCalled();
  });

  it("shows a 429 too_many_sends and stays", async () => {
    const auth = authLoggingIn(async () => {
      throw new ApiError({ status: 429, key: "auth.too_many_sends", params: { limit: 3 } });
    });
    const { router } = renderApp("/entrar", { auth });

    await sendEmail("ana@example.com");

    expect(
      await screen.findByText("Ya enviamos 3 códigos en la última hora. Espera un poco."),
    ).toBeInTheDocument();
    expect(router.state.location.pathname).toBe("/entrar");
  });

  it("goes to Registro from Crea tu cuenta", async () => {
    const { router } = renderApp("/entrar");

    await userEvent.click(await screen.findByRole("link", { name: "¿Primera vez? Crea tu cuenta" }));

    await waitFor(() => expect(router.state.location.pathname).toBe("/registro"));
  });
});
