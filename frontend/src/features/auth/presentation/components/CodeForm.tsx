import { useState } from "react";
import { Button } from "../../../../components/ui/atoms/Button";
import { CodeInput, codeLength } from "../../../../components/ui/molecules/CodeInput";
import { ResendTimer } from "../../../../components/ui/molecules/ResendTimer";
import { t } from "../../../../lib/i18n/translate";

export interface CodeFormProps {
  readonly email: string;
  /** Changes after each failed check: the boxes start over, empty. */
  readonly attempt: number;
  readonly onVerify: (code: string) => void;
  readonly isVerifying: boolean;
  /** Already translated: why the check or the resend failed. */
  readonly error: string | undefined;
  /** The code expired or ran out of attempts: offer a new one right away. */
  readonly asksForNewCode: boolean;
  /** Resolves true when the new code went out. */
  readonly onResend: () => Promise<boolean>;
  /** Changes after each code that went out: the count starts again. */
  readonly sends: number;
}

interface CodeEntryProps {
  readonly onVerify: (code: string) => void;
  readonly isVerifying: boolean;
  readonly isInvalid: boolean;
}

/** The boxes and Verificar; it remembers what was typed so the button can send it. */
function CodeEntry({ onVerify, isVerifying, isInvalid }: CodeEntryProps) {
  const [code, setCode] = useState("");
  return (
    <div className="flex flex-1 flex-col">
      <div className="pt-8">
        <CodeInput onChange={setCode} onComplete={onVerify} disabled={isVerifying} invalid={isInvalid} />
      </div>
      <div className="mt-auto flex flex-col pt-6">
        <Button
          type="button"
          loading={isVerifying}
          disabled={code.length < codeLength}
          onClick={() => onVerify(code)}
        >
          {t("ui.code.verify")}
        </Button>
      </div>
    </div>
  );
}

export function CodeForm({
  email,
  attempt,
  onVerify,
  isVerifying,
  error,
  asksForNewCode,
  onResend,
  sends,
}: CodeFormProps) {
  return (
    <div className="flex flex-1 flex-col px-6 pt-8 md:px-0 md:pt-0">
      <div className="flex flex-col gap-2">
        <h1 className="text-[28px] leading-tight font-semibold tracking-[-0.02em] text-ink">
          {t("ui.code.title")}
        </h1>
        <p className="text-[15px] text-ink-secondary">{t("ui.code.sent_to", { email })}</p>
      </div>
      {error !== undefined && (
        <p role="alert" className="mt-6 rounded-md bg-error-surface p-4 text-sm text-danger-ink">
          {error}
        </p>
      )}
      {asksForNewCode && (
        <Button variant="secondary" className="mt-4" onClick={onResend}>
          {t("ui.code.request_new")}
        </Button>
      )}
      <div className="flex flex-col items-center pt-6">
        <ResendTimer key={sends} onResend={onResend} />
      </div>
      <CodeEntry
        key={attempt}
        onVerify={onVerify}
        isVerifying={isVerifying}
        isInvalid={error !== undefined}
      />
    </div>
  );
}
