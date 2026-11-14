import { AppShell } from "../../../../components/layout/AppShell";
import { BottomNav } from "../../../../components/layout/BottomNav";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { SideNav } from "../../../../components/layout/SideNav";
import { useNavigateTo } from "../../../../lib/router/useNavigateTo";
import { useSession } from "../../application/useSession";
import { firstNameOf } from "../../domain/me";
import { AccessRing } from "../components/AccessRing";
import { WelcomeIntro } from "../components/WelcomeIntro";

/** Bienvenida: the first screen after a first access. It reads only the session. */
export function WelcomeScreen() {
  const session = useSession();
  const navigateTo = useNavigateTo();
  // The route guard already read the session; this only narrows its type.
  if (session.status !== "signedIn") {
    return null;
  }
  const { me } = session;

  return (
    <AppShell
      header={<BrandHeader session={{ name: me.name, onNavigate: navigateTo }} />}
      zone={<AccessRing size={224} />}
      nav={<BottomNav onNavigate={navigateTo} />}
      sideNav={<SideNav account={{ name: me.name, email: me.email }} onNavigate={navigateTo} />}
    >
      <WelcomeIntro firstName={firstNameOf(me.name)} onNavigate={navigateTo} />
    </AppShell>
  );
}
