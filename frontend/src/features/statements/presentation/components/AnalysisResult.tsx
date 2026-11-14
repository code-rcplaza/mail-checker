import type { ReactNode } from "react";
import { sectionPaths } from "../../../../components/layout/navigation";
import { StepIndicator } from "../../../../components/layout/StepIndicator";
import { AppLink, type Navigate } from "../../../../components/ui/atoms/AppLink";
import { CheckIcon, ForwardIcon } from "../../../../components/ui/atoms/icons";
import { EmptyState } from "../../../../components/ui/molecules/EmptyState";
import { SegmentedRing } from "../../../../components/ui/organisms/SegmentedRing";
import { formatMoney, formatMonths } from "../../../../lib/format";
import { type MessageKey, t } from "../../../../lib/i18n/translate";
import type { Found, UploadResult } from "../../domain/upload";

export const reviewPath = "/revision";
export const sourcePath = "/fuente";

/** Checked against the schema of the back, not the handoff: the card's last 4 digits are kept. */
const keptItems: readonly MessageKey[] = [
  "ui.analysis.kept_merchant",
  "ui.analysis.kept_amount",
  "ui.analysis.kept_date",
  "ui.analysis.kept_card",
];
const notKeptItems: readonly MessageKey[] = ["ui.analysis.not_kept_card", "ui.analysis.not_kept_file"];

const linkClasses =
  "flex min-h-[52px] items-center justify-center gap-2 rounded-lg bg-primary px-6 text-[15px] font-semibold text-on-primary";

/** Nothing sure and nothing to confirm: the statement brought no subscriptions. */
export function hasNothingFound(found: Found): boolean {
  return found.sure === 0 && found.possible === 0;
}

/** Análisis, top: step 2 of 3, the ring with the sure and possible ones, and the monthly total. */
export function AnalysisZone({ found }: { readonly found: Found }) {
  const total = found.sure + found.possible;
  const monthly = t(found.monthly.isApproximate ? "ui.analysis.monthly_approximate" : "ui.analysis.monthly", {
    amount: formatMoney(found.monthly),
  });
  const segments = [
    { value: found.sure, series: 1 },
    { value: found.possible, series: 4 },
  ] as const;
  return (
    <div className="flex w-full flex-col gap-6">
      <StepIndicator step={2} total={3} name={t("ui.analysis.step")} />
      <div className="flex items-center gap-4">
        <SegmentedRing variant="normal" segments={segments} label={t("ui.analysis.ring")} size={104}>
          <span className="text-2xl font-semibold text-zone-ink">{total}</span>
        </SegmentedRing>
        <div className="flex min-w-0 flex-1 flex-col gap-1">
          <h1 className="text-[22px] leading-tight font-semibold tracking-[-0.02em] text-zone-ink">
            {t("ui.analysis.title")}
          </h1>
          <p className="text-sm text-hero-ink-secondary">{monthly}</p>
        </div>
      </div>
    </div>
  );
}

function CheckList({ label, items }: { readonly label: string; readonly items: readonly ReactNode[] }) {
  return (
    <div className="flex flex-col gap-2">
      <h3 className="text-sm font-semibold text-ink">{label}</h3>
      <ul aria-label={label} className="flex flex-col gap-2">
        {items.map((item, position) => (
          // biome-ignore lint/suspicious/noArrayIndexKey: fixed texts in a fixed order
          <li key={position} className="flex items-center gap-2 text-sm text-ink-secondary">
            <span className="flex flex-none text-primary">
              <CheckIcon size={16} />
            </span>
            {item}
          </li>
        ))}
      </ul>
    </div>
  );
}

export interface AnalysisResultProps {
  readonly found: Found;
  /** Gone after a reload: then the file part is left out. */
  readonly uploaded: UploadResult | undefined;
  readonly onNavigate: Navigate;
}

export function AnalysisResult({ found, uploaded, onNavigate }: AnalysisResultProps) {
  if (hasNothingFound(found)) {
    return (
      <div className="px-6 pt-8 md:px-0 md:pt-0">
        <EmptyState
          title={t("ui.analysis.empty_title")}
          text={t("ui.analysis.empty_text")}
          action={
            <AppLink to={sourcePath} onNavigate={onNavigate} className={linkClasses}>
              {t("ui.analysis.upload_another")}
            </AppLink>
          }
        />
      </div>
    );
  }
  const done = [
    ...(uploaded === undefined
      ? []
      : [
          t("ui.analysis.read_file", { filename: uploaded.filename }),
          t("ui.analysis.read_months", { months: formatMonths(uploaded.months) }),
        ]),
    t("ui.analysis.looked_for"),
  ];
  const hasDoubts = found.possible > 0;
  return (
    <div className="flex flex-col gap-8 px-6 pt-8 md:px-0 md:pt-0">
      <div className="flex flex-wrap gap-2">
        <span className="inline-flex min-h-9 items-center rounded-full bg-hero-chip px-4 text-sm font-semibold text-zone-ink">
          {t("ui.analysis.sure", { count: found.sure })}
        </span>
        <span className="inline-flex min-h-9 items-center rounded-full bg-alert px-4 text-sm font-semibold text-alert-ink">
          {t("ui.analysis.possible", { count: found.possible })}
        </span>
      </div>
      <section className="flex flex-col gap-4">
        <h2 className="text-[17px] font-semibold text-ink">{t("ui.analysis.done_title")}</h2>
        <CheckList label={t("ui.analysis.done_title")} items={done} />
      </section>
      <section className="grid grid-cols-2 gap-6 border-t border-border pt-8">
        <CheckList label={t("ui.analysis.kept")} items={keptItems.map((key) => t(key))} />
        <CheckList label={t("ui.analysis.not_kept")} items={notKeptItems.map((key) => t(key))} />
      </section>
      <div className="flex flex-col gap-2">
        <AppLink
          to={hasDoubts ? reviewPath : sectionPaths.summary}
          onNavigate={onNavigate}
          className={linkClasses}
        >
          {hasDoubts ? t("ui.analysis.review", { count: found.possible }) : t("ui.analysis.to_summary")}
          <ForwardIcon />
        </AppLink>
        {hasDoubts && found.sure > 0 && (
          <p className="text-center text-xs text-ink-muted">
            {t("ui.analysis.sure_added", { count: found.sure })}
          </p>
        )}
      </div>
    </div>
  );
}
