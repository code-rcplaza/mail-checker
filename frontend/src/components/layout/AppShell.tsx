import type { ReactNode } from "react";
import { classNames } from "../../lib/classNames";

export interface AppShellProps {
  /** BrandHeader: on top of the zone on a phone, above the columns from md:. */
  readonly header: ReactNode;
  /** What the screen shows in the zone (title, ring); from md: it is the hero card on the left. */
  readonly zone?: ReactNode;
  /** BottomNav, below lg:. Access screens have none. */
  readonly nav?: ReactNode;
  /** SideNav, from lg:; it brings the brand, so the header hides there. */
  readonly sideNav?: ReactNode;
  readonly children: ReactNode;
}

/** The two halos of the zone; they only decorate the phone layout. */
function ZoneHalos() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 390 320"
      fill="none"
      className="pointer-events-none absolute top-0 left-0 h-80 w-[390px] stroke-hero-line md:hidden"
    >
      <circle cx="300" cy="40" r="170" strokeWidth="1.5" opacity="0.55" />
      <circle cx="300" cy="40" r="250" strokeWidth="1.5" opacity="0.35" />
    </svg>
  );
}

/**
 * Phone: the zone with its glow and the sheet that rises 24 px over it; the
 * page scrolls and the bottom bar stays fixed above the safe area.
 * From md: (768 px) the zone opens up (`contents`) so the header spans both
 * columns, the hero is a card on the left and the sheet is the right column.
 * From lg: (1024 px) the side bar replaces the bottom one.
 */
export function AppShell({ header, zone, nav, sideNav, children }: AppShellProps) {
  const hasNav = nav !== undefined;
  const hasSideNav = sideNav !== undefined;
  return (
    <div className="flex min-h-dvh">
      {hasSideNav && <div className="hidden lg:block">{sideNav}</div>}
      <div
        data-testid="shell-columns"
        className={classNames(
          "flex min-w-0 flex-1 flex-col",
          "md:mx-auto md:grid md:w-full md:max-w-[1200px] md:grid-cols-2 md:content-start md:items-start md:gap-x-12 md:gap-y-8 md:px-10 md:pt-8 md:pb-16 lg:px-12",
        )}
      >
        <div
          data-testid="shell-zone"
          className="relative overflow-hidden bg-zone px-6 pt-[calc(24px+env(safe-area-inset-top))] pb-12 text-zone-ink md:contents"
        >
          <ZoneHalos />
          <div className="relative flex flex-col items-center gap-6 md:contents">
            <div className={classNames("w-full md:col-span-2", hasSideNav && "lg:hidden")}>{header}</div>
            {zone !== undefined && (
              <section
                data-column="hero"
                className="flex w-full flex-col items-center gap-6 md:rounded-sheet md:bg-zone md:p-10"
              >
                {zone}
              </section>
            )}
          </div>
        </div>
        <main
          data-testid="shell-sheet"
          data-column="content"
          className={classNames(
            "relative z-10 -mt-6 flex flex-1 flex-col rounded-t-sheet bg-page pt-2",
            "md:mt-0 md:rounded-none md:bg-transparent md:pt-0",
            zone === undefined && "md:col-span-2",
            hasNav && "pb-[calc(120px+env(safe-area-inset-bottom))]",
            hasNav && hasSideNav && "lg:pb-0",
          )}
        >
          {children}
        </main>
      </div>
      {hasNav && (
        <div
          className={classNames(
            "fixed inset-x-0 bottom-0 z-20 bg-linear-to-b from-transparent to-page to-[24px] px-4 pt-6 pb-[calc(16px+env(safe-area-inset-bottom))]",
            hasSideNav && "lg:hidden",
          )}
        >
          {nav}
        </div>
      )}
    </div>
  );
}
