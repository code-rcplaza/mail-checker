import { type ReactNode, useId } from "react";

/** What the field inside needs to be tied to its label, help and error. */
export interface FieldControlProps {
  readonly id: string;
  readonly "aria-invalid": true | undefined;
  readonly "aria-describedby": string | undefined;
}

export interface FormFieldProps {
  readonly label: string;
  readonly help?: string | undefined;
  readonly error?: string | undefined;
  readonly children: (control: FieldControlProps) => ReactNode;
}

/** Label, field and messages; with an error, the error comes before the help. */
export function FormField({ label, help, error, children }: FormFieldProps) {
  const id = useId();
  const errorId = `${id}-error`;
  const helpId = `${id}-help`;
  const hasError = error !== undefined && error !== "";
  const hasHelp = help !== undefined && help !== "";
  const describedBy = [hasError && errorId, hasHelp && helpId].filter(Boolean).join(" ");
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="text-sm font-medium text-ink">
        {label}
      </label>
      {children({
        id,
        "aria-invalid": hasError || undefined,
        "aria-describedby": describedBy === "" ? undefined : describedBy,
      })}
      {hasError && (
        <p id={errorId} role="alert" className="text-sm text-danger-ink">
          {error}
        </p>
      )}
      {hasHelp && (
        <p id={helpId} className="text-sm text-ink-muted">
          {help}
        </p>
      )}
    </div>
  );
}
