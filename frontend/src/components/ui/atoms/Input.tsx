import type { InputHTMLAttributes } from "react";
import { classNames } from "../../../lib/classNames";

export type InputProps = InputHTMLAttributes<HTMLInputElement>;

export function Input({ className, ...props }: InputProps) {
  return (
    <input
      {...props}
      className={classNames(
        "min-h-[52px] w-full rounded-md border-[1.5px] border-border-strong bg-field px-4 text-base text-ink",
        "placeholder:text-ink-muted focus-visible:border-primary aria-invalid:border-danger",
        className,
      )}
    />
  );
}
