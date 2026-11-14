import type { ButtonHTMLAttributes } from "react";
import { classNames } from "../../../lib/classNames";
import { t } from "../../../lib/i18n/translate";

export type ButtonVariant = "primary" | "secondary" | "text" | "danger";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  readonly variant?: ButtonVariant;
  /** Working: the button ignores clicks and announces it, keeping its width. */
  readonly loading?: boolean;
}

const variantClasses: Readonly<Record<ButtonVariant, string>> = {
  primary: "bg-primary text-on-primary",
  secondary: "border-[1.5px] border-outline text-outline-ink",
  text: "text-primary",
  danger: "bg-danger text-on-primary",
};

export function Button({
  variant = "primary",
  loading = false,
  disabled,
  className,
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      type="button"
      {...props}
      disabled={disabled === true || loading}
      aria-busy={loading}
      data-variant={variant}
      className={classNames(
        "relative inline-flex min-h-[52px] items-center justify-center gap-2 rounded-lg px-6 text-[15px] font-semibold",
        "transition-opacity disabled:cursor-not-allowed disabled:opacity-60",
        variantClasses[variant],
        className,
      )}
    >
      <span className={classNames("inline-flex items-center gap-2", loading && "invisible")}>{children}</span>
      {loading && (
        <span className="absolute inset-0 flex items-center justify-center">
          <span
            aria-hidden="true"
            className="size-5 rounded-full border-2 border-current border-r-transparent motion-safe:animate-spin"
          />
          <span className="sr-only">{t("ui.loading")}</span>
        </span>
      )}
    </button>
  );
}
