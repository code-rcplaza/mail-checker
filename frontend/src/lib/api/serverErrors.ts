import { t, translateServerKey } from "../i18n/translate";
import { ApiError, NetworkError } from "./client";

/** The i18n key a 400 put on one field, if any. */
export function fieldErrorOf(error: unknown, field: string): string | undefined {
  return error instanceof ApiError ? error.fields[field] : undefined;
}

/** The message for the whole form: the key of the API, or why it did not answer. */
export function serverMessageOf(error: unknown): string {
  if (error instanceof ApiError) {
    return translateServerKey(error.key, error.params);
  }
  if (error instanceof NetworkError) {
    return t("common.unavailable");
  }
  return t("common.internal");
}
