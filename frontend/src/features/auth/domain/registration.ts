import { z } from "zod";

/** The same limit as the back (`auth/domain.NameMaxLength`). */
export const nameMaxLength = 60;

/** The messages are i18n keys: the screen translates them. */
export const registrationSchema = z.object({
  name: z.string().trim().min(1, "auth.name_required").max(nameMaxLength, "auth.name_too_long"),
  email: z.string().trim().pipe(z.email("auth.email_invalid")),
  acceptedTerms: z.boolean().refine((accepted) => accepted, "auth.terms_required"),
});

export type Registration = Readonly<z.infer<typeof registrationSchema>>;

export type RegistrationField = keyof Registration;

export const registrationFields: readonly RegistrationField[] = ["name", "email", "acceptedTerms"];

export const emptyRegistration: Registration = { name: "", email: "", acceptedTerms: false };
