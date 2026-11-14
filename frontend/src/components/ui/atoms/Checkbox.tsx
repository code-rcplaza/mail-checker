import { type InputHTMLAttributes, type ReactNode, useId } from "react";

export interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  readonly label: ReactNode;
  readonly error?: string | undefined;
}

/** A native checkbox inside its label: the label toggles it and space works. */
export function Checkbox({ label, error, ...props }: CheckboxProps) {
  const errorId = useId();
  const hasError = error !== undefined && error !== "";
  return (
    <div className="flex flex-col gap-2">
      <label className="flex min-h-11 items-start gap-4 text-sm leading-normal text-ink-secondary">
        <input
          type="checkbox"
          {...props}
          aria-invalid={hasError || undefined}
          aria-describedby={hasError ? errorId : undefined}
          className="size-6 flex-none accent-primary"
        />
        <span>{label}</span>
      </label>
      {hasError && (
        <p id={errorId} role="alert" className="text-sm text-danger-ink">
          {error}
        </p>
      )}
    </div>
  );
}
