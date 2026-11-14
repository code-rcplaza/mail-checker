import { BackLink } from "../../../../components/layout/BackLink";
import { sectionPaths } from "../../../../components/layout/navigation";
import { StepIndicator } from "../../../../components/layout/StepIndicator";
import type { Navigate } from "../../../../components/ui/atoms/AppLink";
import { SegmentedRing } from "../../../../components/ui/organisms/SegmentedRing";
import { t } from "../../../../lib/i18n/translate";

/** The top of Elegir fuente: back home, step 1 of 3 and the empty ring waiting for charges. */
export function SourceZone({ onNavigate }: { readonly onNavigate: Navigate }) {
  return (
    <div className="flex w-full flex-col gap-6">
      <BackLink to={sectionPaths.summary} label={t("ui.upload.back")} onNavigate={onNavigate} />
      <StepIndicator step={1} total={3} name={t("ui.upload.step")} />
      <div className="flex items-center gap-4">
        <SegmentedRing variant="empty" label={t("ui.access_ring")} size={104} />
        <div className="flex min-w-0 flex-1 flex-col gap-1">
          <h1 className="text-[22px] leading-tight font-semibold tracking-[-0.02em] text-zone-ink">
            {t("ui.upload.title")}
          </h1>
          <p className="text-sm text-hero-ink-secondary">{t("ui.upload.subtitle")}</p>
        </div>
      </div>
    </div>
  );
}
