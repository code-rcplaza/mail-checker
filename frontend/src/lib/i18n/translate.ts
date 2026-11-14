import { es } from "./es";

export type MessageKey = keyof typeof es;

export type MessageParams = Readonly<Record<string, string | number>>;

function isMessageKey(key: string): key is MessageKey {
  return Object.hasOwn(es, key);
}

function fill(text: string, params: MessageParams): string {
  return text.replace(/\{(\w+)\}/g, (placeholder, name: string) => String(params[name] ?? placeholder));
}

/** Translates a key the front knows at compile time. */
export function t(key: MessageKey, params: MessageParams = {}): string {
  return fill(es[key], params);
}

/**
 * Translates a key that came from the back. One the catalog does not have yet
 * shows the generic error, never the raw key.
 */
export function translateServerKey(key: string, params: Readonly<Record<string, unknown>> = {}): string {
  if (!isMessageKey(key)) {
    return t("common.internal");
  }
  const printable = Object.fromEntries(
    Object.entries(params).filter((entry): entry is [string, string | number] =>
      ["string", "number"].includes(typeof entry[1]),
    ),
  );
  return t(key, printable);
}
