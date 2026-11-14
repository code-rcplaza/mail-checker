import { useEffect, useState } from "react";
import { t } from "../../../lib/i18n/translate";
import { Button } from "../atoms/Button";

export const resendWaitSeconds = 30;

function clockOf(seconds: number): string {
  return `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, "0")}`;
}

export interface ResendTimerProps {
  /**
   * Plain callback: the count restarts at once. A promise: the button waits
   * for it and the count restarts only if it resolves `true` (the code went out).
   */
  readonly onResend: () => unknown;
}

/** The resend button waits 30 s after each send, counting down. */
export function ResendTimer({ onResend }: ResendTimerProps) {
  const [secondsLeft, setSecondsLeft] = useState(resendWaitSeconds);
  const [isSending, setIsSending] = useState(false);

  useEffect(() => {
    if (secondsLeft === 0) {
      return;
    }
    const timer = window.setTimeout(() => setSecondsLeft(secondsLeft - 1), 1000);
    return () => window.clearTimeout(timer);
  }, [secondsLeft]);

  async function resend() {
    const result = onResend();
    if (!(result instanceof Promise)) {
      setSecondsLeft(resendWaitSeconds);
      return;
    }
    setIsSending(true);
    const wentOut: unknown = await result;
    setIsSending(false);
    if (wentOut === true) {
      setSecondsLeft(resendWaitSeconds);
    }
  }

  const isWaiting = secondsLeft > 0;
  return (
    <Button variant="text" disabled={isWaiting || isSending} onClick={resend}>
      {isWaiting ? t("ui.code.resend_in", { time: clockOf(secondsLeft) }) : t("ui.code.resend")}
    </Button>
  );
}
