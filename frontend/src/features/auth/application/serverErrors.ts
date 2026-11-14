import { ApiError, NetworkError } from "../../../lib/api/client";
import { t, translateServerKey } from "../../../lib/i18n/translate";

/** The field errors of a 400 that belong to the form, as `[field, key]`. */
export function serverFieldErrors<Field extends string>(
  error: unknown,
  fields: readonly Field[],
): ReadonlyArray<readonly [Field, string]> {
  if (!(error instanceof ApiError)) {
    return [];
  }
  return fields.flatMap((field) => {
    const key = error.fields[field];
    return key === undefined ? [] : [[field, key] as const];
  });
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
