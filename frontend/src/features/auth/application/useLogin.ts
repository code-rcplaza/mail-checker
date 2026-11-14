import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import { translateServerKey } from "../../../lib/i18n/translate";
import { emptyLogin, type Login, loginSchema } from "../domain/login";
import { serverFieldErrors, serverMessageOf } from "./serverErrors";

export interface UseLoginOptions {
  /** Called with the email once the code is on its way. */
  readonly onCodeSent: (email: string) => void;
}

const loginFields = ["email"] as const;

/** Login: validates the email, asks for a code once and reports what the server rejected. */
export function useLogin({ onCodeSent }: UseLoginOptions) {
  const { auth } = useGateways();
  const form = useForm({ resolver: zodResolver(loginSchema), defaultValues: emptyLogin });
  const [formError, setFormError] = useState<string | undefined>(undefined);

  function showServerError(error: unknown) {
    const fieldErrors = serverFieldErrors(error, loginFields);
    for (const [field, key] of fieldErrors) {
      form.setError(field, { message: key });
    }
    setFormError(fieldErrors.length > 0 ? undefined : serverMessageOf(error));
  }

  const sending = useMutation({
    mutationFn: ({ email }: Login) => auth.login(email),
    onSuccess: (_answer, { email }) => onCodeSent(email),
    onError: showServerError,
  });

  function send(login: Login) {
    if (sending.isPending) {
      return;
    }
    setFormError(undefined);
    sending.mutate(login);
  }

  const emailError = form.formState.errors.email?.message;
  return {
    fields: { email: form.register("email") },
    errors: emailError === undefined ? {} : { email: translateServerKey(emailError) },
    formError,
    isSubmitting: sending.isPending,
    submit: form.handleSubmit(send),
  };
}
