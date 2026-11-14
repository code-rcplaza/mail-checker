import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { ApiError } from "../../../lib/api/client";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import { translateServerKey } from "../../../lib/i18n/translate";
import { asksForNewCode, type CodeEntry } from "../domain/code";
import { serverFieldErrors, serverMessageOf } from "./serverErrors";

export interface UseVerifyCodeOptions {
  readonly email: string;
  /** Called once the session is open. */
  readonly onVerified: () => void;
}

interface VerifyFailure {
  readonly message: string;
  readonly asksForNewCode: boolean;
}

/** A field error (the code or the email) says more than "check the fields". */
function failureOf(error: unknown): VerifyFailure {
  const [fieldError] = serverFieldErrors(error, ["code", "email"]);
  return {
    message: fieldError === undefined ? serverMessageOf(error) : translateServerKey(fieldError[1]),
    asksForNewCode: error instanceof ApiError && asksForNewCode(error.key),
  };
}

/** Código: checks the six digits once; after a failure the boxes start over. */
export function useVerifyCode({ email, onVerified }: UseVerifyCodeOptions) {
  const { auth } = useGateways();
  const [failure, setFailure] = useState<VerifyFailure | undefined>(undefined);
  const [attempt, setAttempt] = useState(0);

  const verifying = useMutation({
    mutationFn: (entry: CodeEntry) => auth.verifyCode(entry),
    onSuccess: onVerified,
    onError: (error) => {
      setFailure(failureOf(error));
      setAttempt((previous) => previous + 1);
    },
  });

  function verify(code: string) {
    if (verifying.isPending) {
      return;
    }
    setFailure(undefined);
    verifying.mutate({ email, code });
  }

  return {
    verify,
    isVerifying: verifying.isPending,
    error: failure?.message,
    asksForNewCode: failure?.asksForNewCode ?? false,
    /** Changes after each failure: the boxes are drawn again, empty. */
    attempt,
    clearError: () => setFailure(undefined),
  };
}
