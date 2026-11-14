import { act, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import { renderApp } from "../../../../app/renderApp";
import { ApiError } from "../../../../lib/api/client";
import { gatewaysWith } from "../../../../lib/gateways/renderWithGateways";
import type { AuthGateway } from "../../application/ports";

function authWith(overrides: Partial<AuthGateway>): AuthGateway {
  return {
    ...gatewaysWith({}).auth,
    verifyCode: vi.fn(async () => {}),
    resendCode: vi.fn(async () => {}),
    ...overrides,
  };
}

const box = (n: number) => screen.getByRole("textbox", { name: `Dígito ${n} de 6` });

async function typeCode(code: string) {
  await userEvent.click(await screen.findByRole("textbox", { name: "Dígito 1 de 6" }));
  await userEvent.keyboard(code);
}

describe("Code screen", () => {
  afterEach(() => vi.useRealTimers());

  it("says where the code went", async () => {
    renderApp("/codigo?email=ana%40example.com&flow=login", { auth: authWith({}) });

    expect(await screen.findByRole("heading", { name: "Revisa tu correo" })).toBeInTheDocument();
    expect(screen.getByText("Enviamos un código de 6 dígitos a ana@example.com.")).toBeInTheDocument();
  });

  it.each([
    ["register", "/bienvenida"],
    ["login", "/resumen"],
  ])("sends the six digits once and, coming from %s, goes to %s", async (flow, destination) => {
    const auth = authWith({});
    const { router } = renderApp(`/codigo?email=ana%40example.com&flow=${flow}`, { auth });

    await typeCode("123456");

    await waitFor(() => expect(router.state.location.pathname).toBe(destination));
    expect(auth.verifyCode).toHaveBeenCalledOnce();
    expect(auth.verifyCode).toHaveBeenCalledWith({ email: "ana@example.com", code: "123456" });
  });

  it("does not send again while the first check is on its way", async () => {
    const auth = authWith({ verifyCode: vi.fn(() => new Promise<void>(() => {})) });
    renderApp("/codigo?email=ana%40example.com&flow=login", { auth });

    await typeCode("123456");
    await userEvent.click(screen.getByRole("button", { name: /Verificar/ }));

    expect(auth.verifyCode).toHaveBeenCalledOnce();
  });

  it("shows a wrong code and empties the boxes", async () => {
    const auth = authWith({
      verifyCode: vi.fn(async () => {
        throw new ApiError({ status: 400, key: "auth.code_invalid" });
      }),
    });
    const { router } = renderApp("/codigo?email=ana%40example.com&flow=login", { auth });

    await typeCode("123456");

    expect(await screen.findByText("El código no es correcto.")).toBeInTheDocument();
    for (const n of [1, 2, 3, 4, 5, 6]) {
      expect(box(n)).toHaveValue("");
    }
    expect(router.state.location.pathname).toBe("/codigo");
  });

  it.each([
    ["auth.code_expired", "El código venció. Pide uno nuevo."],
    ["auth.code_attempts_exceeded", "Usaste todos los intentos. Pide un código nuevo."],
  ])("offers a new code after %s", async (key, message) => {
    const auth = authWith({
      verifyCode: vi.fn(async () => {
        throw new ApiError({ status: 400, key });
      }),
    });
    renderApp("/codigo?email=ana%40example.com&flow=login", { auth });

    await typeCode("123456");
    expect(await screen.findByText(message)).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Pedir otro código" }));

    expect(auth.resendCode).toHaveBeenCalledWith("ana@example.com");
  });

  it("goes back to Login when it was opened without an email", async () => {
    const { router } = renderApp("/codigo");

    await waitFor(() => expect(router.state.location.pathname).toBe("/entrar"));
  });

  it("resends once the timer is done, says so and restarts the count", async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    const auth = authWith({});
    renderApp("/codigo?email=ana%40example.com&flow=login", { auth });
    await screen.findByRole("button", { name: "Reenviar en 0:30" });

    for (let second = 0; second < 30; second++) {
      await act(async () => vi.advanceTimersByTime(1000));
    }
    await userEvent.click(screen.getByRole("button", { name: "Reenviar código" }));

    expect(auth.resendCode).toHaveBeenCalledWith("ana@example.com");
    expect(await screen.findByRole("status")).toHaveTextContent("Te enviamos un código nuevo.");
    expect(screen.getByRole("button", { name: "Reenviar en 0:30" })).toBeDisabled();
  });

  it("shows a 429 on resend and does not restart the count", async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    const auth = authWith({
      resendCode: vi.fn(async () => {
        throw new ApiError({ status: 429, key: "auth.too_many_sends", params: { limit: 3 } });
      }),
    });
    renderApp("/codigo?email=ana%40example.com&flow=login", { auth });
    await screen.findByRole("button", { name: "Reenviar en 0:30" });

    for (let second = 0; second < 30; second++) {
      await act(async () => vi.advanceTimersByTime(1000));
    }
    await userEvent.click(screen.getByRole("button", { name: "Reenviar código" }));

    expect(
      await screen.findByText("Ya enviamos 3 códigos en la última hora. Espera un poco."),
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Reenviar código" })).toBeEnabled();
  });
});
