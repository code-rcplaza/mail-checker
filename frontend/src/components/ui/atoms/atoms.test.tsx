import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { Badge } from "./Badge";
import { Button } from "./Button";
import { Checkbox } from "./Checkbox";
import { initialOf, ServiceAvatar, seriesClassOf } from "./ServiceAvatar";

describe("Button", () => {
  it.each([
    ["primary", "bg-primary"],
    ["secondary", "border-outline"],
    ["text", "text-primary"],
    ["danger", "bg-danger"],
  ] as const)("styles the %s variant", (variant, expectedClass) => {
    render(<Button variant={variant}>Seguir</Button>);

    expect(screen.getByRole("button", { name: "Seguir" })).toHaveClass(expectedClass);
  });

  it("is at least 44 px tall to touch", () => {
    render(<Button>Seguir</Button>);

    expect(screen.getByRole("button")).toHaveClass("min-h-[52px]");
  });

  it.each([
    ["disabled", { disabled: true }],
    ["loading", { loading: true }],
  ])("ignores clicks while %s", async (_state, props) => {
    const onClick = vi.fn();
    render(
      <Button {...props} onClick={onClick}>
        Seguir
      </Button>,
    );

    await userEvent.click(screen.getByRole("button"));

    expect(onClick).not.toHaveBeenCalled();
  });

  it("announces loading and keeps the label for its width", () => {
    render(<Button loading>Seguir</Button>);

    const button = screen.getByRole("button");
    expect(button).toHaveAttribute("aria-busy", "true");
    expect(button).toHaveTextContent("Cargando");
    expect(screen.getByText("Seguir")).toHaveClass("invisible");
  });
});

describe("Checkbox", () => {
  it("toggles from its label and with space", async () => {
    render(<Checkbox label="Acepto los términos" />);
    const box = screen.getByRole("checkbox", { name: "Acepto los términos" });

    await userEvent.click(screen.getByText("Acepto los términos"));
    expect(box).toBeChecked();

    await userEvent.keyboard(" ");
    expect(box).not.toBeChecked();
  });

  it("ties its error to the box", () => {
    render(<Checkbox label="Acepto" error="Para seguir, acepta los términos." />);

    const box = screen.getByRole("checkbox");
    expect(box).toBeInvalid();
    expect(box).toHaveAccessibleDescription("Para seguir, acepta los términos.");
  });
});

describe("Badge", () => {
  it.each([
    ["toCancel", "Por cancelar"],
    ["rise", "Alza"],
    ["possible", "Posible"],
  ] as const)("says %s in words", (kind, text) => {
    render(<Badge kind={kind} />);

    expect(screen.getByText(text)).toBeInTheDocument();
  });
});

describe("ServiceAvatar", () => {
  it("shows the initial", () => {
    render(<ServiceAvatar name="Netflix" />);

    expect(screen.getByText("N")).toBeInTheDocument();
  });

  it("gives the same name the same color", () => {
    expect(seriesClassOf("Netflix")).toBe(seriesClassOf(" netflix "));
    expect(seriesClassOf("Netflix")).toMatch(/^bg-series-[1-6]$/);
  });

  it("shows a neutral symbol for an empty name", () => {
    expect(initialOf("")).toBe("•");
    expect(initialOf("  *  ")).toBe("•");
    render(<ServiceAvatar name="" />);
    expect(screen.getByText("•")).toBeInTheDocument();
  });
});
