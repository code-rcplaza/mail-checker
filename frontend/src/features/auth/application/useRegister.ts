import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { type FieldErrors, useForm } from "react-hook-form";
import { fieldErrorOf, serverMessageOf } from "../../../lib/api/serverErrors";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import { translateServerKey } from "../../../lib/i18n/translate";
import {
  emptyRegistration,
  nameMaxLength,
  type Registration,
  type RegistrationField,
  registrationFields,
  registrationSchema,
} from "../domain/registration";

export interface UseRegisterOptions {
  /** Called with the email once the code is on its way. */
  readonly onRegistered: (email: string) => void;
}

/** Field messages are keys (zod or the back); `{max}` fills auth.name_too_long. */
function translateFieldErrors(errors: FieldErrors<Registration>): Partial<Record<RegistrationField, string>> {
  return Object.fromEntries(
    registrationFields.flatMap((field) => {
      const key = errors[field]?.message;
      return key === undefined ? [] : [[field, translateServerKey(key, { max: nameMaxLength })]];
    }),
  );
}

/** Registro: validates with zod, sends once and reports what the server rejected. */
export function useRegister({ onRegistered }: UseRegisterOptions) {
  const { auth } = useGateways();
  const form = useForm({ resolver: zodResolver(registrationSchema), defaultValues: emptyRegistration });
  const [formError, setFormError] = useState<string | undefined>(undefined);

  function showServerError(error: unknown) {
    const rejected = registrationFields.flatMap((field) => {
      const key = fieldErrorOf(error, field);
      return key === undefined ? [] : [{ field, key }];
    });
    for (const { field, key } of rejected) {
      form.setError(field, { message: key });
    }
    setFormError(rejected.length > 0 ? undefined : serverMessageOf(error));
  }

  const sending = useMutation({
    mutationFn: (registration: Registration) => auth.register(registration),
    onSuccess: (_answer, registration) => onRegistered(registration.email),
    onError: showServerError,
  });

  function send(registration: Registration) {
    if (sending.isPending) {
      return;
    }
    setFormError(undefined);
    sending.mutate(registration);
  }

  return {
    fields: {
      name: form.register("name"),
      email: form.register("email"),
      acceptedTerms: form.register("acceptedTerms"),
    },
    errors: translateFieldErrors(form.formState.errors),
    formError,
    isSubmitting: sending.isPending,
    submit: form.handleSubmit(send),
  };
}
