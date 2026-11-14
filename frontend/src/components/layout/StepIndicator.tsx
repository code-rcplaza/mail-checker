import { classNames } from "../../lib/classNames";
import { t } from "../../lib/i18n/translate";

export interface StepIndicatorProps {
  readonly step: number;
  readonly total: number;
  /** The name of the current step ("Análisis"). */
  readonly name: string;
}

/** A whole number inside [min, max]; anything else (NaN included) becomes min. */
export function clampWhole(value: number, range: { readonly min: number; readonly max: number }): number {
  if (!Number.isFinite(value)) {
    return range.min;
  }
  return Math.min(range.max, Math.max(range.min, Math.trunc(value)));
}

/** The flows have 3 steps; a bigger total would only draw bars nobody reads. */
const maxSteps = 10;

/** A step out of range shows the nearest one instead of breaking. */
export function StepIndicator({ step, total, name }: StepIndicatorProps) {
  const steps = clampWhole(total, { min: 1, max: maxSteps });
  const current = clampWhole(step, { min: 1, max: steps });
  const bars = Array.from({ length: steps }, (_, index) => ({
    position: index + 1,
    isDone: index < current,
  }));
  return (
    <div className="flex w-full flex-col gap-2">
      <p className="text-xs font-semibold tracking-[0.08em] text-hero-ink-secondary uppercase">
        {t("ui.step", { step: current, total: steps, name })}
      </p>
      <div aria-hidden="true" className="flex gap-2">
        {bars.map(({ position, isDone }) => (
          <span
            key={position}
            data-testid="step-bar"
            data-done={isDone}
            className={classNames("h-1 flex-1 rounded-full", isDone ? "bg-zone-ink" : "bg-hero-line")}
          />
        ))}
      </div>
    </div>
  );
}
