import { AppShell } from "../../../../components/layout/AppShell";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { sectionPaths } from "../../../../components/layout/navigation";
import { Toast } from "../../../../components/ui/molecules/Toast";
import { useNavigateTo } from "../../../../lib/router/useNavigateTo";
import { useThemeToggle } from "../../../../lib/theme/useThemeToggle";
import { useResendCode } from "../../application/useResendCode";
import { useVerifyCode } from "../../application/useVerifyCode";
import type { AccessFlow } from "../../domain/code";
import { AccessRing } from "../components/AccessRing";
import { CodeForm } from "../components/CodeForm";

/** Bienvenida greets a first access; everyone else goes to their Resumen. */
export const welcomePath = "/bienvenida";

export function pathAfterVerifying(flow: AccessFlow | undefined): string {
  return flow === "register" ? welcomePath : sectionPaths.summary;
}

export interface CodeScreenProps {
  readonly email: string;
  readonly flow: AccessFlow | undefined;
}

/** Código: checks the code, resends it and opens the session. */
export function CodeScreen({ email, flow }: CodeScreenProps) {
  const navigateTo = useNavigateTo();
  const theme = useThemeToggle();
  const verification = useVerifyCode({ email, onVerified: () => navigateTo(pathAfterVerifying(flow)) });
  const resending = useResendCode({ email });

  async function resendCode(): Promise<boolean> {
    const wentOut = await resending.resend();
    if (wentOut) {
      verification.clearError();
    }
    return wentOut;
  }

  return (
    <AppShell header={<BrandHeader theme={theme} />} zone={<AccessRing />}>
      <CodeForm
        email={email}
        attempt={verification.attempt}
        onVerify={verification.verify}
        isVerifying={verification.isVerifying}
        error={resending.error ?? verification.error}
        asksForNewCode={verification.asksForNewCode}
        onResend={resendCode}
        sends={resending.notice?.id ?? 0}
      />
      {resending.notice !== undefined && (
        <Toast
          key={resending.notice.id}
          message={resending.notice.message}
          onClose={resending.dismissNotice}
        />
      )}
    </AppShell>
  );
}
