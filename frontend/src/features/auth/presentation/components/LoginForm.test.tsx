import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { LoginForm, type LoginFormProps } from "./LoginForm";

function props(overrides: Partial<LoginFormProps> = {}): LoginFormProps {
  return {
    fields: { email: { name: "email", onChange: async () => {}, onBlur: async () => {}, ref: () => {} } },
    errors: {},
    formError: undefined,
    isSubmitting: false,
    onSubmit: (event) => event.preventDefault(),
    onNavigate: () => {},
    ...overrides,
  };
}

describe("LoginForm", () => {
  it("asks only for the email", () => {
    render(<LoginForm {...props()} />);

    expect(screen.getByRole("heading", { name: "Entra a Fokko" })).toBeInTheDocument();
    expect(screen.getByLabelText("Tu correo")).toHaveAttribute("autocomplete", "email");
    expect(screen.getAllByRole("textbox")).toHaveLength(1);
    expect(screen.getByRole("button", { name: "Continuar" })).toBeEnabled();
  });

  it("shows the field error and the general one", () => {
    render(
      <LoginForm {...props({ errors: { email: "Escribe un correo válido." }, formError: "Demasiados." })} />,
    );

    expect(screen.getByLabelText("Tu correo")).toHaveAccessibleDescription("Escribe un correo válido.");
    expect(screen.getAllByRole("alert").map((alert) => alert.textContent)).toContain("Demasiados.");
  });

  it("blocks the button while sending", () => {
    render(<LoginForm {...props({ isSubmitting: true })} />);

    expect(screen.getByRole("button", { name: /Continuar/ })).toBeDisabled();
  });

  it("takes first-timers to Registro", async () => {
    const onNavigate = vi.fn();
    render(<LoginForm {...props({ onNavigate })} />);

    await userEvent.click(screen.getByRole("link", { name: "¿Primera vez? Crea tu cuenta" }));

    expect(onNavigate).toHaveBeenCalledWith("/registro");
  });
});
