import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { fieldErrorOf, serverMessageOf } from "../../../lib/api/serverErrors";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import { translateServerKey } from "../../../lib/i18n/translate";
import { emptyLogin, type Login, loginSchema } from "../domain/login";

export interface UseLoginOptions {
  /** Called with the email once the code is on its way. */
  readonly onCodeSent: (email: string) => void;
}

/** Login: validates the email, asks for a code once and reports what the server rejected. */
export function useLogin({ onCodeSent }: UseLoginOptions) {
  const { auth } = useGateways();
  const form = useForm({ resolver: zodResolver(loginSchema), defaultValues: emptyLogin });
  const [formError, setFormError] = useState<string | undefined>(undefined);

  function showServerError(error: unknown) {
    const emailKey = fieldErrorOf(error, "email");
    if (emailKey !== undefined) {
      form.setError("email", { message: emailKey });
    }
    setFormError(emailKey === undefined ? serverMessageOf(error) : undefined);
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
