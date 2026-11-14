import { type ClipboardEvent, type KeyboardEvent, useRef, useState } from "react";
import { t } from "../../../lib/i18n/translate";

export const codeLength = 6;

const emptyDigits = (): string[] => Array.from({ length: codeLength }, () => "");

function digitsOnly(text: string): string {
  return text.replace(/\D/g, "");
}

export interface CodeInputProps {
  /** Called once the six digits are in. */
  readonly onComplete: (code: string) => void;
  /** Called on every change, with the digits typed so far. */
  readonly onChange?: ((code: string) => void) | undefined;
  readonly disabled?: boolean;
  readonly invalid?: boolean;
}

/** Six boxes for the login code: digits only, it advances, pastes and goes back. */
export function CodeInput({ onComplete, onChange, disabled = false, invalid = false }: CodeInputProps) {
  const [digits, setDigits] = useState(emptyDigits);
  const boxes = useRef<Array<HTMLInputElement | null>>([]);

  function focusBox(index: number) {
    boxes.current[Math.min(Math.max(index, 0), codeLength - 1)]?.focus();
  }

  function update(next: string[]) {
    setDigits(next);
    onChange?.(next.join(""));
    if (next.every((digit) => digit !== "")) {
      onComplete(next.join(""));
    }
  }

  function typeAt(index: number, typed: string) {
    const digit = digitsOnly(typed).slice(-1);
    if (typed === "") {
      update(digits.map((current, position) => (position === index ? "" : current)));
      return;
    }
    if (digit === "") {
      return;
    }
    update(digits.map((current, position) => (position === index ? digit : current)));
    focusBox(index + 1);
  }

  function pasteFrom(index: number, event: ClipboardEvent<HTMLInputElement>) {
    event.preventDefault();
    const pasted = digitsOnly(event.clipboardData.getData("text")).slice(0, codeLength - index);
    if (pasted === "") {
      return;
    }
    update(digits.map((current, position) => pasted[position - index] ?? current));
    focusBox(index + pasted.length);
  }

  function goBackFrom(index: number, event: KeyboardEvent<HTMLInputElement>) {
    if (event.key !== "Backspace" || digits[index] !== "") {
      return;
    }
    event.preventDefault();
    update(digits.map((current, position) => (position === index - 1 ? "" : current)));
    focusBox(index - 1);
  }

  return (
    <div className="flex justify-between gap-2">
      {digits.map((digit, index) => (
        <input
          // biome-ignore lint/suspicious/noArrayIndexKey: the six boxes are fixed positions
          key={index}
          ref={(element) => {
            boxes.current[index] = element;
          }}
          value={digit}
          disabled={disabled}
          inputMode="numeric"
          autoComplete={index === 0 ? "one-time-code" : "off"}
          maxLength={2}
          aria-label={t("ui.code.digit", { n: index + 1, total: codeLength })}
          aria-invalid={invalid || undefined}
          onChange={(event) => typeAt(index, event.target.value)}
          onPaste={(event) => pasteFrom(index, event)}
          onKeyDown={(event) => goBackFrom(index, event)}
          className="size-12 rounded-md border-[1.5px] border-border-strong bg-field text-center text-xl font-semibold text-ink focus-visible:border-primary aria-invalid:border-danger"
        />
      ))}
    </div>
  );
}
