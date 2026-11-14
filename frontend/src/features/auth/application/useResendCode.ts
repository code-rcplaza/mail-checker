import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { serverMessageOf } from "../../../lib/api/serverErrors";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import { t } from "../../../lib/i18n/translate";

export interface ResendNotice {
  /** A new id per notice, so a new toast replaces the old one and restarts. */
  readonly id: number;
  readonly message: string;
}

/** "Reenviar código": tells whether it went out, so the count restarts only then. */
export function useResendCode({ email }: { readonly email: string }) {
  const { auth } = useGateways();
  const [error, setError] = useState<string | undefined>(undefined);
  const [notice, setNotice] = useState<ResendNotice | undefined>(undefined);
  const resending = useMutation({ mutationFn: () => auth.resendCode(email) });

  async function resend(): Promise<boolean> {
    setError(undefined);
    try {
      await resending.mutateAsync();
    } catch (caught) {
      setError(serverMessageOf(caught));
      return false;
    }
    setNotice((previous) => ({ id: (previous?.id ?? 0) + 1, message: t("ui.code.resent") }));
    return true;
  }

  return { resend, error, notice, dismissNotice: () => setNotice(undefined) };
}
