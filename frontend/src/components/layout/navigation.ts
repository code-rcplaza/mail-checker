import type { MessageKey } from "../../lib/i18n/translate";
import { type IconSize, SubsIcon, SummaryIcon } from "../ui/atoms/icons";

/** The sections of the app with a session and where each one lives. */
export const sectionPaths = {
  summary: "/resumen",
  subs: "/subs",
  profile: "/perfil",
} as const;

export type NavSection = "summary" | "subs";

export type { Navigate } from "../ui/atoms/AppLink";

export interface NavEntry {
  readonly section: NavSection;
  readonly label: MessageKey;
  readonly Icon: (props: { readonly size?: IconSize }) => React.JSX.Element;
}

/** MVP: only Resumen and Subs (Calendario and Fuentes come later). Both bars use it. */
export const navEntries: readonly NavEntry[] = [
  { section: "summary", label: "ui.nav.summary", Icon: SummaryIcon },
  { section: "subs", label: "ui.nav.subs", Icon: SubsIcon },
];
