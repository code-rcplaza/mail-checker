import { classNames } from "../../lib/classNames";
import { t } from "../../lib/i18n/translate";
import { AppLink } from "./AppLink";
import { BrandMark, initialsOf } from "./BrandHeader";
import { type Navigate, type NavSection, navEntries, sectionPaths } from "./navigation";

export interface SideNavAccount {
  readonly name: string;
  readonly email: string;
}

export interface SideNavProps {
  readonly current: NavSection;
  readonly account: SideNavAccount;
  readonly onNavigate: Navigate;
}

/** Perfil is not a section: the footer with the account leads there. */
function AccountFooter({ account, onNavigate }: Omit<SideNavProps, "current">) {
  return (
    <AppLink
      to={sectionPaths.profile}
      onNavigate={onNavigate}
      aria-label={t("ui.nav.profile")}
      className="-ml-2 mt-auto flex items-center gap-4 rounded-lg p-2"
    >
      <span className="flex size-10 flex-none items-center justify-center rounded-full bg-avatar text-sm font-semibold text-avatar-ink">
        {initialsOf(account.name)}
      </span>
      <span className="flex min-w-0 flex-col gap-1">
        <span className="truncate text-sm font-semibold">{account.name}</span>
        <span className="truncate text-xs text-ink-muted">{account.email}</span>
      </span>
    </AppLink>
  );
}

/** Desktop bar (from lg:): 256 px, sticky, with a border on the right. */
export function SideNav({ current, account, onNavigate }: SideNavProps) {
  return (
    <aside className="sticky top-0 flex h-dvh w-64 flex-none flex-col gap-8 border-r border-border px-6 py-8 text-ink">
      <BrandMark />
      <nav aria-label={t("ui.nav.label")} className="flex flex-col gap-1">
        {navEntries.map(({ section, label, Icon }) => {
          const isCurrent = section === current;
          return (
            <AppLink
              key={section}
              to={sectionPaths[section]}
              onNavigate={onNavigate}
              aria-current={isCurrent ? "page" : undefined}
              className={classNames(
                "flex min-h-11 items-center gap-4 rounded-lg px-4 text-[15px] font-semibold",
                isCurrent ? "bg-hero-line text-zone-ink" : "text-ink-secondary",
              )}
            >
              <Icon size={20} />
              {t(label)}
            </AppLink>
          );
        })}
      </nav>
      <AccountFooter account={account} onNavigate={onNavigate} />
    </aside>
  );
}
