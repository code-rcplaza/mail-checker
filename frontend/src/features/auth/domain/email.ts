import { z } from "zod";

/** The email of Registro and Login: trimmed, and its message is an i18n key. */
export const emailField = z.string().trim().pipe(z.email("auth.email_invalid"));
