import { classNames } from "../../lib/classNames";
import { t } from "../../lib/i18n/translate";
import { AppLink } from "./AppLink";
import { type Navigate, type NavSection, navEntries, sectionPaths } from "./navigation";

export interface BottomNavProps {
  readonly current: NavSection;
  readonly onNavigate: Navigate;
}

/** The current one shows its name; the others only their icon (named for readers). */
export function BottomNav({ current, onNavigate }: BottomNavProps) {
  return (
    <nav
      aria-label={t("ui.nav.label")}
      className="flex items-center justify-around gap-1 rounded-full border border-border bg-nav p-2"
    >
      {navEntries.map(({ section, label, Icon }) => {
        const isCurrent = section === current;
        return (
          <AppLink
            key={section}
            to={sectionPaths[section]}
            onNavigate={onNavigate}
            aria-label={t(label)}
            aria-current={isCurrent ? "page" : undefined}
            className={classNames(
              "flex min-h-12 items-center justify-center gap-2 rounded-full text-sm font-semibold",
              isCurrent ? "bg-primary px-4 text-on-primary" : "w-12 text-ink-muted",
            )}
          >
            <Icon size={20} />
            {isCurrent && <span aria-hidden="true">{t(label)}</span>}
          </AppLink>
        );
      })}
    </nav>
  );
}
