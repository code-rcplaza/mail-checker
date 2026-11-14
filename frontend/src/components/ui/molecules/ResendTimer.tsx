import { useEffect, useState } from "react";
import { t } from "../../../lib/i18n/translate";
import { Button } from "../atoms/Button";

export const resendWaitSeconds = 30;

function clockOf(seconds: number): string {
  return `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, "0")}`;
}

/** The resend button waits 30 s after each send, counting down. */
export function ResendTimer({ onResend }: { readonly onResend: () => void }) {
  const [secondsLeft, setSecondsLeft] = useState(resendWaitSeconds);

  useEffect(() => {
    if (secondsLeft === 0) {
      return;
    }
    const timer = window.setTimeout(() => setSecondsLeft(secondsLeft - 1), 1000);
    return () => window.clearTimeout(timer);
  }, [secondsLeft]);

  function resend() {
    onResend();
    setSecondsLeft(resendWaitSeconds);
  }

  const isWaiting = secondsLeft > 0;
  return (
    <Button variant="text" disabled={isWaiting} onClick={resend}>
      {isWaiting ? t("ui.code.resend_in", { time: clockOf(secondsLeft) }) : t("ui.code.resend")}
    </Button>
  );
}
