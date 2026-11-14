import { t } from "../../lib/i18n/translate";
import { AppLink } from "../ui/atoms/AppLink";
import { MoonIcon, SunIcon } from "../ui/atoms/icons";
import { type Navigate, sectionPaths } from "./navigation";

const neutralInitials = "•";

/** First letters of the first and last word: "Carla de la Vidal" → "CV". */
export function initialsOf(name: string): string {
  const words = name.trim().split(/\s+/).filter(Boolean);
  const first = words[0];
  const last = words.length > 1 ? words[words.length - 1] : undefined;
  if (first === undefined) {
    return neutralInitials;
  }
  return [first, last]
    .map((word) => word?.charAt(0) ?? "")
    .join("")
    .toLocaleUpperCase("es-CL");
}

export interface BrandHeaderSession {
  readonly name: string;
  readonly onNavigate: Navigate;
}

export interface BrandHeaderTheme {
  readonly isDark: boolean;
  readonly onToggle: () => void;
}

/** With a session the avatar shows; without one, the theme button (access screens). */
export type BrandHeaderProps =
  | { readonly session: BrandHeaderSession; readonly theme?: never }
  | { readonly theme: BrandHeaderTheme; readonly session?: never };

/** Wordmark and product pill; they take the color of where they sit. */
export function BrandMark() {
  return (
    <div className="flex min-w-0 items-start gap-2">
      <span className="font-wordmark text-[28px] leading-none font-semibold">{t("ui.brand.name")}</span>
      <span className="inline-flex min-h-5 items-center rounded-full border-[1.5px] border-current px-2 text-[10px] leading-none font-semibold tracking-[0.08em] whitespace-nowrap uppercase">
        {t("ui.brand.product")}
      </span>
    </div>
  );
}

function ProfileAvatar({ session }: { readonly session: BrandHeaderSession }) {
  return (
    <AppLink
      to={sectionPaths.profile}
      onNavigate={session.onNavigate}
      aria-label={t("ui.nav.profile")}
      className="flex size-11 flex-none items-center justify-center rounded-full bg-hero-line text-sm font-semibold text-zone-ink"
    >
      <span aria-hidden="true">{initialsOf(session.name)}</span>
    </AppLink>
  );
}

function ThemeButton({ theme }: { readonly theme: BrandHeaderTheme }) {
  return (
    <button
      type="button"
      onClick={theme.onToggle}
      aria-label={t(theme.isDark ? "ui.theme.to_light" : "ui.theme.to_dark")}
      className="flex size-11 flex-none items-center justify-center rounded-full text-zone-ink"
    >
      {theme.isDark ? <SunIcon /> : <MoonIcon />}
    </button>
  );
}

export function BrandHeader({ session, theme }: BrandHeaderProps) {
  return (
    <header className="flex w-full items-center justify-between gap-2">
      <BrandMark />
      {session !== undefined && <ProfileAvatar session={session} />}
      {theme !== undefined && <ThemeButton theme={theme} />}
    </header>
  );
}
