import { StepIndicator } from "../../../../components/layout/StepIndicator";
import { Skeleton } from "../../../../components/ui/atoms/Skeleton";
import { SegmentedRing } from "../../../../components/ui/organisms/SegmentedRing";
import { t } from "../../../../lib/i18n/translate";

/** Cargando, top: step 2 of 3 and the spinning ring (still when motion is reduced). */
export function UploadLoadingZone() {
  return (
    <div className="flex w-full flex-col gap-6">
      <StepIndicator step={2} total={3} name={t("ui.analysis.step")} />
      <div className="flex items-center gap-4">
        <SegmentedRing variant="loading" size={104}>
          <span className="text-xs font-semibold text-hero-ink-secondary">
            {t("ui.loading_statement.ring")}
          </span>
        </SegmentedRing>
        <div className="flex min-w-0 flex-1 flex-col gap-1">
          <h1 className="text-[22px] leading-tight font-semibold tracking-[-0.02em] text-zone-ink">
            {t("ui.loading_statement.title")}
          </h1>
          <p className="text-sm text-hero-ink-secondary">{t("ui.loading_statement.text")}</p>
        </div>
      </div>
    </div>
  );
}

const placeholderRows = ["first", "second", "third"] as const;

/** Cargando, sheet: rows that stand for the charges about to appear. */
export function UploadLoadingContent() {
  return (
    <div aria-busy="true" className="flex flex-col px-6 pt-8 md:px-0 md:pt-0">
      <p role="status" className="sr-only">
        {t("ui.loading_statement.status")}
      </p>
      {placeholderRows.map((row) => (
        <div key={row} className="flex items-center gap-4 py-4">
          <Skeleton className="size-10 flex-none rounded-full" />
          <div className="flex flex-1 flex-col gap-2">
            <Skeleton className="h-3.5 w-[55%]" />
            <Skeleton className="h-3 w-[30%]" />
          </div>
        </div>
      ))}
    </div>
  );
}
