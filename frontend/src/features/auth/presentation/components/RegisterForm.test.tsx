import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { RegisterForm, type RegisterFormProps } from "./RegisterForm";

function binding(name: string) {
  return { name, onChange: async () => {}, onBlur: async () => {}, ref: () => {} };
}

function props(overrides: Partial<RegisterFormProps> = {}): RegisterFormProps {
  return {
    fields: { name: binding("name"), email: binding("email"), acceptedTerms: binding("acceptedTerms") },
    errors: {},
    formError: undefined,
    isSubmitting: false,
    onSubmit: (event) => event.preventDefault(),
    onNavigate: () => {},
    ...overrides,
  };
}

describe("RegisterForm", () => {
  it("asks for name, email and the terms", () => {
    render(<RegisterForm {...props()} />);

    expect(screen.getByRole("heading", { name: "Crea tu cuenta" })).toBeInTheDocument();
    expect(screen.getByLabelText("Tu nombre")).toHaveAttribute("autocomplete", "name");
    expect(screen.getByLabelText("Tu correo")).toHaveAttribute("type", "email");
    expect(
      screen.getByRole("checkbox", { name: "Acepto los términos y la política de privacidad." }),
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Crear cuenta" })).toBeEnabled();
  });

  it("shows each error on its field and the general one", () => {
    render(
      <RegisterForm
        {...props({
          errors: { name: "Escribe tu nombre.", acceptedTerms: "Para seguir, acepta los términos." },
          formError: "Demasiados intentos.",
        })}
      />,
    );

    expect(screen.getByLabelText("Tu nombre")).toHaveAccessibleDescription("Escribe tu nombre.");
    expect(screen.getByRole("checkbox")).toHaveAttribute("aria-invalid", "true");
    expect(screen.getAllByRole("alert").map((alert) => alert.textContent)).toContain("Demasiados intentos.");
  });

  it("blocks the button while sending", () => {
    render(<RegisterForm {...props({ isSubmitting: true })} />);

    expect(screen.getByRole("button", { name: /Crear cuenta/ })).toBeDisabled();
  });

  it("takes the people who already have an account to Login", async () => {
    const onNavigate = vi.fn();
    render(<RegisterForm {...props({ onNavigate })} />);

    await userEvent.click(screen.getByRole("link", { name: "¿Ya tienes cuenta? Entra" }));

    expect(onNavigate).toHaveBeenCalledWith("/entrar");
  });
});
