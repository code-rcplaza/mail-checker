import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { renderApp } from "../../../../app/renderApp";
import { ApiError } from "../../../../lib/api/client";
import { gatewaysWith } from "../../../../lib/gateways/renderWithGateways";
import type { AuthGateway } from "../../application/ports";

function authRegistering(register: AuthGateway["register"]): AuthGateway {
  return { ...gatewaysWith({}).auth, register: vi.fn(register) };
}

async function fillValidForm() {
  await userEvent.type(await screen.findByLabelText("Tu nombre"), "Ana Prueba");
  await userEvent.type(screen.getByLabelText("Tu correo"), "ana@example.com");
  await userEvent.click(screen.getByRole("checkbox"));
}

describe("Register screen", () => {
  it("sends valid data and goes to Código with the email", async () => {
    const auth = authRegistering(async () => {});
    const { router } = renderApp("/registro", { auth });

    await fillValidForm();
    await userEvent.click(screen.getByRole("button", { name: "Crear cuenta" }));

    await waitFor(() => expect(router.state.location.pathname).toBe("/codigo"));
    expect(router.state.location.search).toEqual({ email: "ana@example.com" });
    expect(auth.register).toHaveBeenCalledWith({
      name: "Ana Prueba",
      email: "ana@example.com",
      acceptedTerms: true,
    });
  });

  it("checks the fields before calling the server", async () => {
    const auth = authRegistering(async () => {});
    renderApp("/registro", { auth });

    await userEvent.click(await screen.findByRole("button", { name: "Crear cuenta" }));

    expect(await screen.findByText("Escribe tu nombre.")).toBeInTheDocument();
    expect(screen.getByText("Escribe un correo válido.")).toBeInTheDocument();
    expect(screen.getByText("Para seguir, acepta los términos.")).toBeInTheDocument();
    expect(auth.register).not.toHaveBeenCalled();
  });

  it("blocks a second click while the first one is on its way", async () => {
    const auth = authRegistering(() => new Promise(() => {}));
    renderApp("/registro", { auth });

    await fillValidForm();
    await userEvent.click(screen.getByRole("button", { name: "Crear cuenta" }));

    await waitFor(() => expect(screen.getByRole("button", { name: /Crear cuenta/ })).toBeDisabled());
    await userEvent.click(screen.getByRole("button", { name: /Crear cuenta/ }));
    expect(auth.register).toHaveBeenCalledOnce();
  });

  it("shows a 429 too_many_sends and stays", async () => {
    const auth = authRegistering(async () => {
      throw new ApiError({ status: 429, key: "auth.too_many_sends", params: { limit: 3 } });
    });
    const { router } = renderApp("/registro", { auth });

    await fillValidForm();
    await userEvent.click(screen.getByRole("button", { name: "Crear cuenta" }));

    expect(
      await screen.findByText("Ya enviamos 3 códigos en la última hora. Espera un poco."),
    ).toBeInTheDocument();
    expect(router.state.location.pathname).toBe("/registro");
  });

  it("puts each server field error on its field", async () => {
    const auth = authRegistering(async () => {
      throw new ApiError({
        status: 400,
        key: "common.validation_failed",
        fields: { email: "auth.email_invalid", name: "auth.name_too_long" },
      });
    });
    renderApp("/registro", { auth });

    await fillValidForm();
    await userEvent.click(screen.getByRole("button", { name: "Crear cuenta" }));

    await waitFor(() =>
      expect(screen.getByLabelText("Tu correo")).toHaveAccessibleDescription(/Escribe un correo válido./),
    );
    expect(screen.getByLabelText("Tu nombre")).toHaveAccessibleDescription(
      "El nombre puede tener hasta 60 caracteres.",
    );
  });

  it("offers the theme button, not the avatar", async () => {
    renderApp("/registro");

    expect(await screen.findByRole("button", { name: /Cambiar a tema/ })).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: "Perfil" })).not.toBeInTheDocument();
  });
});
