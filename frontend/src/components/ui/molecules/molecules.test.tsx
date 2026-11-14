import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { Input } from "../atoms/Input";
import { Card } from "./Card";
import { FormField } from "./FormField";
import { ListRow } from "./ListRow";

describe("FormField", () => {
  it("ties the label to the field", () => {
    render(<FormField label="Correo">{(control) => <Input {...control} />}</FormField>);

    expect(screen.getByLabelText("Correo")).toBeInstanceOf(HTMLInputElement);
  });

  it("marks the field invalid, announces the error and puts it before the help", () => {
    render(
      <FormField label="Correo" help="Te mandaremos un código." error="Escribe un correo válido.">
        {(control) => <Input {...control} />}
      </FormField>,
    );

    const field = screen.getByLabelText("Correo");
    expect(field).toBeInvalid();
    expect(screen.getByRole("alert")).toHaveTextContent("Escribe un correo válido.");
    expect(field).toHaveAccessibleDescription("Escribe un correo válido. Te mandaremos un código.");
    const error = screen.getByText("Escribe un correo válido.");
    const help = screen.getByText("Te mandaremos un código.");
    expect(error.compareDocumentPosition(help) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
  });

  it("is valid without an error", () => {
    render(<FormField label="Correo">{(control) => <Input {...control} />}</FormField>);

    expect(screen.getByLabelText("Correo")).toBeValid();
  });
});

describe("Card", () => {
  it.each([
    ["hero", "bg-hero-chip"],
    ["alert", "bg-alert"],
    ["tip", "bg-field"],
  ] as const)("styles the %s variant", (variant, expectedClass) => {
    render(<Card variant={variant}>Contenido</Card>);

    expect(screen.getByText("Contenido")).toHaveClass(expectedClass);
  });
});

describe("ListRow", () => {
  it("runs its action on click and on Enter", async () => {
    const onSelect = vi.fn();
    render(<ListRow title="Netflix" onSelect={onSelect} />);

    await userEvent.click(screen.getByRole("button", { name: /Netflix/ }));
    screen.getByRole("button").focus();
    await userEvent.keyboard("{Enter}");

    expect(onSelect).toHaveBeenCalledTimes(2);
  });

  it("reads title, subtitle and value in that order", () => {
    render(<ListRow title="Netflix" subtitle="BCI · crédito" value="$7.990" />);

    expect(screen.getByText("Netflix").closest("div")).toHaveTextContent(/^NetflixBCI · crédito\$7\.990$/);
  });

  it("is at least 44 px tall to touch", () => {
    render(<ListRow title="Netflix" onSelect={() => {}} />);

    expect(screen.getByRole("button")).toHaveClass("min-h-11");
  });

  it("is not a button without an action", () => {
    render(<ListRow title="Netflix" />);

    expect(screen.queryByRole("button")).toBeNull();
  });
});
