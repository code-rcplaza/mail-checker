import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { firstNameOf } from "../../domain/me";
import { WelcomeIntro } from "./WelcomeIntro";

describe("WelcomeIntro", () => {
  it("greets with the name", () => {
    render(<WelcomeIntro firstName="Ana" onNavigate={() => {}} />);

    expect(screen.getByText("Hola, Ana")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Tu mes, a la vista" })).toBeInTheDocument();
  });

  it("takes Empezar to Elegir fuente", async () => {
    const onNavigate = vi.fn();
    render(<WelcomeIntro firstName="Ana" onNavigate={onNavigate} />);

    await userEvent.click(screen.getByRole("link", { name: "Empezar" }));

    expect(onNavigate).toHaveBeenCalledWith("/fuente");
  });
});

describe("firstNameOf", () => {
  it.each([
    ["Ana Pérez", "Ana"],
    ["  Ana  ", "Ana"],
    ["", ""],
  ])("takes the first name of %j", (name, first) => {
    expect(firstNameOf(name)).toBe(first);
  });
});
