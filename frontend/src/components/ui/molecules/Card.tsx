import type { ReactNode } from "react";
import { classNames } from "../../../lib/classNames";

export type CardVariant = "hero" | "alert" | "tip";

const variantClasses: Readonly<Record<CardVariant, string>> = {
  hero: "bg-hero-chip text-hero-ink border-hero-line",
  alert: "bg-alert text-alert-ink border-alert-border",
  tip: "bg-field text-ink-secondary border-border",
};

/** Only for the hero, an alert or a tip: the rest of the screens are not boxed. */
export function Card({ variant, children }: { readonly variant: CardVariant; readonly children: ReactNode }) {
  return (
    <section data-variant={variant} className={classNames("rounded-lg border p-4", variantClasses[variant])}>
      {children}
    </section>
  );
}
