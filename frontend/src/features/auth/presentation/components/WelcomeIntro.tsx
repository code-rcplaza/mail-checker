import { AppLink, type Navigate } from "../../../../components/ui/atoms/AppLink";
import { CheckIcon, ForwardIcon } from "../../../../components/ui/atoms/icons";
import { type MessageKey, t } from "../../../../lib/i18n/translate";

/** Elegir fuente (FOKKO-0099-SUBS): in the MVP, only uploading a statement. */
export const sourcePath = "/fuente";

/** The handoff talks about emails; the MVP reads statements (`docs/mvp.md`). */
const points: readonly MessageKey[] = [
  "ui.welcome.point_statement",
  "ui.welcome.point_rises",
  "ui.welcome.point_privacy",
];

export interface WelcomeIntroProps {
  readonly firstName: string;
  readonly onNavigate: Navigate;
}

export function WelcomeIntro({ firstName, onNavigate }: WelcomeIntroProps) {
  return (
    <div className="flex flex-1 flex-col px-6 pt-8 md:px-0 md:pt-0">
      <div className="flex flex-col gap-4">
        <p className="text-sm font-semibold text-accent">{t("ui.welcome.greeting", { name: firstName })}</p>
        <h1 className="text-4xl leading-tight font-semibold tracking-[-0.02em] text-ink md:text-[44px]">
          {t("ui.welcome.title")}
        </h1>
        <p className="text-base text-ink-secondary">{t("ui.welcome.text")}</p>
      </div>
      <ul className="flex flex-col gap-4 pt-6">
        {points.map((point) => (
          <li key={point} className="flex items-center gap-4 text-[15px] text-ink">
            <span className="flex flex-none text-primary">
              <CheckIcon />
            </span>
            {t(point)}
          </li>
        ))}
      </ul>
      <div className="mt-auto flex flex-col pt-6">
        <AppLink
          to={sourcePath}
          onNavigate={onNavigate}
          className="flex min-h-[52px] items-center justify-center gap-2 rounded-lg bg-primary px-6 text-[15px] font-semibold text-on-primary"
        >
          {t("ui.welcome.start")}
          <ForwardIcon />
        </AppLink>
      </div>
    </div>
  );
}
