import { act, fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { Skeleton } from "../atoms/Skeleton";
import { CodeInput } from "./CodeInput";
import { EmptyState } from "./EmptyState";
import { ResendTimer } from "./ResendTimer";
import { SegmentedControl } from "./SegmentedControl";
import { Toast, toastDuration } from "./Toast";

function TwoNotices() {
  const [notice, setNotice] = useState({ id: 1, message: "Primero" });
  return (
    <>
      <button type="button" onClick={() => setNotice({ id: 2, message: "Segundo" })}>
        otro
      </button>
      <Toast key={notice.id} message={notice.message} onClose={() => {}} />
    </>
  );
}

const options = [
  { value: "light", label: "Claro" },
  { value: "dark", label: "Oscuro" },
  { value: "system", label: "Sistema" },
];

function Appearance({ onChange }: { readonly onChange: (value: string) => void }) {
  const [value, setValue] = useState("light");
  return (
    <SegmentedControl
      label="Apariencia"
      options={options}
      value={value}
      onChange={(next) => {
        setValue(next);
        onChange(next);
      }}
    />
  );
}

describe("Toast", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  it("is announced and closes by itself", () => {
    const onClose = vi.fn();
    render(<Toast message="Guardado" onClose={onClose} />);

    expect(screen.getByRole("status")).toHaveTextContent("Guardado");
    act(() => vi.advanceTimersByTime(toastDuration));
    expect(onClose).toHaveBeenCalledOnce();
  });

  it("closes with its button and with Esc", () => {
    const onClose = vi.fn();
    render(<Toast message="Guardado" onClose={onClose} />);

    fireEvent.click(screen.getByRole("button", { name: "Cerrar" }));
    fireEvent.keyDown(document, { key: "Escape" });

    expect(onClose).toHaveBeenCalledTimes(2);
  });

  it("lets a new notice replace the previous one", () => {
    render(<TwoNotices />);

    fireEvent.click(screen.getByRole("button", { name: "otro" }));

    expect(screen.getAllByRole("status")).toHaveLength(1);
    expect(screen.getByRole("status")).toHaveTextContent("Segundo");
  });
});

describe("Skeleton", () => {
  it("is hidden from screen readers and only pulses with motion allowed", () => {
    const { container } = render(<Skeleton className="h-4" />);

    const skeleton = container.firstElementChild;
    expect(skeleton).toHaveAttribute("aria-hidden", "true");
    expect(skeleton).toHaveClass("motion-safe:animate-pulse");
    expect(skeleton).not.toHaveClass("animate-pulse");
  });
});

describe("EmptyState", () => {
  it("shows title, text and action", () => {
    render(
      <EmptyState
        title="Sin suscripciones"
        text="Sube una cartola."
        action={<button type="button">Subir</button>}
      />,
    );

    expect(screen.getByRole("heading", { name: "Sin suscripciones" })).toBeInTheDocument();
    expect(screen.getByText("Sube una cartola.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Subir" })).toBeInTheDocument();
  });
});

describe("SegmentedControl", () => {
  it("calls the handler with the chosen value and marks it", async () => {
    const onChange = vi.fn();
    render(<Appearance onChange={onChange} />);

    await userEvent.click(screen.getByRole("radio", { name: "Sistema" }));

    expect(onChange).toHaveBeenCalledWith("system");
    expect(screen.getByRole("radio", { name: "Sistema" })).toBeChecked();
    expect(screen.getByRole("radio", { name: "Claro" })).not.toBeChecked();
  });

  it("moves between options with the arrows", async () => {
    const onChange = vi.fn();
    render(<Appearance onChange={onChange} />);
    screen.getByRole("radio", { name: "Claro" }).focus();

    await userEvent.keyboard("{ArrowRight}");
    expect(screen.getByRole("radio", { name: "Oscuro" })).toHaveFocus();
    await userEvent.keyboard("{ArrowLeft}{ArrowLeft}");

    expect(onChange).toHaveBeenLastCalledWith("system");
    expect(screen.getByRole("radio", { name: "Sistema" })).toBeChecked();
  });
});

describe("CodeInput", () => {
  const box = (n: number) => screen.getByRole("textbox", { name: `Dígito ${n} de 6` });

  it("moves to the next box after a digit and ignores letters", async () => {
    render(<CodeInput onComplete={() => {}} />);

    await userEvent.click(box(1));
    await userEvent.keyboard("4a");

    expect(box(1)).toHaveValue("4");
    expect(box(2)).toHaveValue("");
    expect(box(2)).toHaveFocus();
  });

  it("fills the six boxes from a paste and completes", async () => {
    const onComplete = vi.fn();
    render(<CodeInput onComplete={onComplete} />);

    await userEvent.click(box(1));
    await userEvent.paste("123456");

    for (const [index, digit] of [..."123456"].entries()) {
      expect(box(index + 1)).toHaveValue(digit);
    }
    expect(onComplete).toHaveBeenCalledWith("123456");
  });

  it("goes back to the previous box on Backspace in an empty one", async () => {
    render(<CodeInput onComplete={() => {}} />);
    await userEvent.click(box(1));
    await userEvent.keyboard("7");

    await userEvent.keyboard("{Backspace}");

    expect(box(1)).toHaveFocus();
    expect(box(1)).toHaveValue("");
  });

  it("clears a filled box on Backspace", async () => {
    render(<CodeInput onComplete={() => {}} />);
    await userEvent.click(box(1));
    await userEvent.keyboard("12");

    await userEvent.click(box(1));
    await userEvent.keyboard("{Backspace}");

    expect(box(1)).toHaveValue("");
    expect(box(2)).toHaveValue("2");
  });

  it("tells every change of the code, complete or not", async () => {
    const onChange = vi.fn();
    render(<CodeInput onComplete={() => {}} onChange={onChange} />);

    await userEvent.click(box(1));
    await userEvent.keyboard("12{Backspace}");

    expect(onChange.mock.calls.map(([code]) => code)).toEqual(["1", "12", "1"]);
  });

  it("calls onComplete with the code once the six digits are typed", async () => {
    const onComplete = vi.fn();
    render(<CodeInput onComplete={onComplete} />);

    await userEvent.click(box(1));
    await userEvent.keyboard("90210");
    expect(onComplete).not.toHaveBeenCalled();
    await userEvent.keyboard("8");

    expect(onComplete).toHaveBeenCalledWith("902108");
  });
});

describe("ResendTimer", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  function advance(seconds: number) {
    for (let i = 0; i < seconds; i++) {
      act(() => vi.advanceTimersByTime(1000));
    }
  }

  it("counts down with the button disabled, then enables it", () => {
    render(<ResendTimer onResend={() => {}} />);
    expect(screen.getByRole("button", { name: "Reenviar en 0:30" })).toBeDisabled();

    advance(29);
    expect(screen.getByRole("button", { name: "Reenviar en 0:01" })).toBeDisabled();
    advance(1);

    expect(screen.getByRole("button", { name: "Reenviar código" })).toBeEnabled();
  });

  it("waits for the result of an async resend and restarts only when it went out", async () => {
    let finish: (sent: boolean) => void = () => {};
    const onResend = vi.fn(() => new Promise<boolean>((resolve) => (finish = resolve)));
    render(<ResendTimer onResend={onResend} />);
    advance(30);

    fireEvent.click(screen.getByRole("button", { name: "Reenviar código" }));
    expect(screen.getByRole("button", { name: "Reenviar código" })).toBeDisabled();
    await act(async () => finish(false));

    expect(screen.getByRole("button", { name: "Reenviar código" })).toBeEnabled();
    fireEvent.click(screen.getByRole("button", { name: "Reenviar código" }));
    await act(async () => finish(true));
    expect(screen.getByRole("button", { name: "Reenviar en 0:30" })).toBeDisabled();
  });

  it("restarts the countdown after a resend", () => {
    const onResend = vi.fn();
    render(<ResendTimer onResend={onResend} />);
    advance(30);

    fireEvent.click(screen.getByRole("button", { name: "Reenviar código" }));

    expect(onResend).toHaveBeenCalledOnce();
    expect(screen.getByRole("button", { name: "Reenviar en 0:30" })).toBeDisabled();
  });
});
