import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { AppShell } from "./AppShell";
import { BackLink } from "./BackLink";
import { BottomNav } from "./BottomNav";
import { BrandHeader, initialsOf } from "./BrandHeader";
import { sectionPaths } from "./navigation";
import { SideNav } from "./SideNav";
import { StepIndicator } from "./StepIndicator";

describe("AppShell", () => {
  it("puts the sheet with 28 px corners over the top zone", () => {
    render(
      <AppShell header={<p>marca</p>} zone={<p>hero</p>} nav={<p>barra</p>}>
        <p>contenido</p>
      </AppShell>,
    );

    const zone = screen.getByTestId("shell-zone");
    const sheet = screen.getByTestId("shell-sheet");
    expect(zone).toHaveTextContent("marcahero");
    expect(zone).toHaveClass("bg-zone");
    expect(sheet).toHaveTextContent("contenido");
    expect(sheet).toHaveClass("rounded-t-sheet", "-mt-6", "z-10");
  });

  it("keeps the bottom navigation fixed while the content scrolls", () => {
    render(
      <AppShell header={<p>marca</p>} nav={<p>barra</p>}>
        <p>contenido</p>
      </AppShell>,
    );

    const nav = screen.getByText("barra").parentElement;
    expect(nav).toHaveClass("fixed", "bottom-0");
    expect(screen.getByTestId("shell-sheet")).toHaveClass("pb-[calc(120px+env(safe-area-inset-bottom))]");
  });

  it("respects the safe areas of a notch", () => {
    render(
      <AppShell header={<p>marca</p>} nav={<p>barra</p>}>
        <p>contenido</p>
      </AppShell>,
    );

    expect(screen.getByTestId("shell-zone")).toHaveClass("pt-[calc(24px+env(safe-area-inset-top))]");
    expect(screen.getByText("barra").parentElement).toHaveClass(
      "pb-[calc(16px+env(safe-area-inset-bottom))]",
    );
  });

  it("leaves no room for a navigation it does not have", () => {
    render(
      <AppShell header={<p>marca</p>}>
        <p>contenido</p>
      </AppShell>,
    );

    expect(screen.getByTestId("shell-sheet")).not.toHaveClass("pb-[calc(120px+env(safe-area-inset-bottom))]");
  });
});

describe("BrandHeader", () => {
  it("shows the wordmark, the product pill and the initials with a session", () => {
    render(<BrandHeader session={{ name: "Carla Vidal", onNavigate: () => {} }} />);

    expect(screen.getByText("Fokko")).toBeInTheDocument();
    expect(screen.getByText("Subs tracker")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Perfil" })).toHaveTextContent("CV");
    expect(screen.queryByRole("button")).not.toBeInTheDocument();
  });

  it("shows the theme button instead of the avatar without a session", async () => {
    const onToggleTheme = vi.fn();
    render(<BrandHeader theme={{ isDark: false, onToggle: onToggleTheme }} />);

    await userEvent.click(screen.getByRole("button", { name: "Cambiar a tema oscuro" }));

    expect(onToggleTheme).toHaveBeenCalledOnce();
    expect(screen.queryByRole("link", { name: "Perfil" })).not.toBeInTheDocument();
  });

  it("names the theme button for the theme it switches to", () => {
    render(<BrandHeader theme={{ isDark: true, onToggle: () => {} }} />);

    expect(screen.getByRole("button", { name: "Cambiar a tema claro" })).toBeInTheDocument();
  });

  it("takes the avatar to Perfil", async () => {
    const onNavigate = vi.fn();
    render(<BrandHeader session={{ name: "Carla Vidal", onNavigate }} />);

    await userEvent.click(screen.getByRole("link", { name: "Perfil" }));

    expect(onNavigate).toHaveBeenCalledWith(sectionPaths.profile);
  });

  it.each([
    ["Carla Vidal", "CV"],
    ["  carla   de la vidal ", "CV"],
    ["Carla", "C"],
    ["Ñandú", "Ñ"],
    ["   ", "•"],
  ])("takes the initials of %j", (name, initials) => {
    expect(initialsOf(name)).toBe(initials);
  });
});

describe("BottomNav", () => {
  it("marks the current section", () => {
    render(<BottomNav current="summary" onNavigate={() => {}} />);

    expect(screen.getByRole("link", { name: "Resumen" })).toHaveAttribute("aria-current", "page");
    expect(screen.getByRole("link", { name: "Subs" })).not.toHaveAttribute("aria-current");
  });

  it("goes to the other section on a tap", async () => {
    const onNavigate = vi.fn();
    render(<BottomNav current="summary" onNavigate={onNavigate} />);

    await userEvent.click(screen.getByRole("link", { name: "Subs" }));

    expect(onNavigate).toHaveBeenCalledWith(sectionPaths.subs);
  });

  it("gives every item a name and at least 44 px", () => {
    render(<BottomNav current="subs" onNavigate={() => {}} />);

    for (const link of screen.getAllByRole("link")) {
      expect(link).toHaveAccessibleName();
      expect(link).toHaveClass("min-h-12");
    }
  });

  it("leaves a click with a modifier to the browser", () => {
    const onNavigate = vi.fn();
    render(<BottomNav current="summary" onNavigate={onNavigate} />);

    fireEvent.click(screen.getByRole("link", { name: "Subs" }), { metaKey: true });

    expect(onNavigate).not.toHaveBeenCalled();
    expect(screen.getByRole("link", { name: "Subs" })).toHaveAttribute("href", sectionPaths.subs);
  });
});

describe("BackLink", () => {
  it("goes to its destination", async () => {
    const onNavigate = vi.fn();
    render(<BackLink to="/fuente" label="Fuente" onNavigate={onNavigate} />);

    await userEvent.click(screen.getByRole("link", { name: "Fuente" }));

    expect(onNavigate).toHaveBeenCalledWith("/fuente");
    expect(screen.getByRole("link", { name: "Fuente" })).toHaveClass("min-h-11");
  });
});

describe("StepIndicator", () => {
  it("says which step it is", () => {
    render(<StepIndicator step={2} total={3} name="Análisis" />);

    expect(screen.getByText("Paso 2 de 3 · Análisis")).toBeInTheDocument();
    expect(screen.getAllByTestId("step-bar").map((bar) => bar.dataset.done)).toEqual([
      "true",
      "true",
      "false",
    ]);
  });

  it.each([
    [5, 3, "Paso 3 de 3 · Revisión", 3],
    [-1, 3, "Paso 1 de 3 · Revisión", 3],
    [1, 0, "Paso 1 de 1 · Revisión", 1],
    [Number.NaN, 2, "Paso 1 de 2 · Revisión", 2],
    [1, 1e9, "Paso 1 de 10 · Revisión", 10],
  ])("does not break with step %d of %d", (step, total, text, bars) => {
    render(<StepIndicator step={step} total={total} name="Revisión" />);

    expect(screen.getByText(text)).toBeInTheDocument();
    expect(screen.getAllByTestId("step-bar")).toHaveLength(bars);
  });
});

describe("SideNav", () => {
  const account = { name: "Carla Vidal", email: "carla@example.com" };

  it("brings the brand, Resumen, Subs and the footer with the account", () => {
    render(<SideNav current="summary" account={account} onNavigate={() => {}} />);

    expect(screen.getByText("Fokko")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Resumen" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Subs" })).toBeInTheDocument();
    const footer = screen.getByRole("link", { name: "Perfil" });
    expect(footer).toHaveTextContent("CV");
    expect(footer).toHaveTextContent("Carla Vidal");
    expect(footer).toHaveTextContent("carla@example.com");
  });

  it("takes the footer to Perfil", async () => {
    const onNavigate = vi.fn();
    render(<SideNav current="summary" account={account} onNavigate={onNavigate} />);

    await userEvent.click(screen.getByRole("link", { name: "Perfil" }));

    expect(onNavigate).toHaveBeenCalledWith(sectionPaths.profile);
  });

  it("marks the current section", () => {
    render(<SideNav current="subs" account={account} onNavigate={() => {}} />);

    expect(screen.getByRole("link", { name: "Subs" })).toHaveAttribute("aria-current", "page");
    expect(screen.getByRole("link", { name: "Resumen" })).not.toHaveAttribute("aria-current");
  });
});

// jsdom does not evaluate media queries: these check the breakpoints the
// classes ask for; the look at 390, 768 and 1280 px is reviewed by eye.
describe("AppShell breakpoints", () => {
  function renderSignedIn() {
    return render(
      <AppShell header={<p>marca</p>} zone={<p>hero</p>} nav={<p>barra</p>} sideNav={<p>lateral</p>}>
        <p>contenido</p>
      </AppShell>,
    );
  }

  it("shows BottomNav below lg: and SideNav from lg:", () => {
    renderSignedIn();

    expect(screen.getByText("lateral").parentElement).toHaveClass("hidden", "lg:block");
    expect(screen.getByText("barra").parentElement).toHaveClass("lg:hidden");
    expect(screen.getByText("marca").parentElement).toHaveClass("lg:hidden");
  });

  it("uses the two columns of the desktop from 768 px, with the bottom bar still there", () => {
    renderSignedIn();

    const columns = screen.getByTestId("shell-columns");
    expect(columns).toHaveClass("md:grid", "md:grid-cols-2", "md:gap-x-12", "md:max-w-[1200px]");
    expect(screen.getByText("barra").parentElement).not.toHaveClass("md:hidden");
  });

  it("puts the hero on the left and the screen on the right", () => {
    renderSignedIn();

    const columns = [...screen.getByTestId("shell-columns").querySelectorAll("[data-column]")];
    expect(columns.map((column) => [column.getAttribute("data-column"), column.textContent])).toEqual([
      ["hero", "hero"],
      ["content", "contenido"],
    ]);
    expect(screen.getByTestId("shell-zone")).toHaveClass("md:contents");
  });

  it("stacks the columns on a phone", () => {
    renderSignedIn();

    expect(screen.getByTestId("shell-columns")).toHaveClass("flex", "flex-col");
    expect(screen.getByTestId("shell-columns")).not.toHaveClass("grid");
  });

  it("shows neither side nor bottom bar on the access screens", () => {
    render(
      <AppShell header={<p>marca</p>} zone={<p>hero</p>}>
        <p>formulario</p>
      </AppShell>,
    );

    expect(screen.queryByRole("navigation")).not.toBeInTheDocument();
    expect(screen.getByText("marca").parentElement).not.toHaveClass("lg:hidden");
    expect(screen.getByTestId("shell-columns")).toHaveClass("md:grid-cols-2");
  });
});
