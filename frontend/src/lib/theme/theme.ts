export const themePreferences = ["light", "dark", "system"] as const;

export type ThemePreference = (typeof themePreferences)[number];

const storageKey = "fokko.theme";

/** Where the preference lives; the browser's localStorage in the app. */
export interface PreferenceStorage {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
}

function isThemePreference(value: string | null): value is ThemePreference {
  return themePreferences.some((preference) => preference === value);
}

/** Reads the stored preference; blocked or broken storage means "system". */
export function readThemePreference(storage: PreferenceStorage): ThemePreference {
  try {
    const stored = storage.getItem(storageKey);
    return isThemePreference(stored) ? stored : "system";
  } catch {
    return "system";
  }
}

/** Stores the preference; when storage is blocked or full it lasts only until reload. */
export function storeThemePreference(storage: PreferenceStorage, preference: ThemePreference): boolean {
  try {
    storage.setItem(storageKey, preference);
    return true;
  } catch {
    return false;
  }
}

/**
 * Applies the preference to the document: "light" and "dark" win over the
 * system; "system" leaves it to `prefers-color-scheme` in the tokens.
 */
export function applyThemePreference(root: HTMLElement, preference: ThemePreference): void {
  if (preference === "system") {
    root.removeAttribute("data-theme");
    return;
  }
  root.setAttribute("data-theme", preference);
}

/** The browser's storage, reached on each call: even touching it can throw when blocked. */
export const browserStorage: PreferenceStorage = {
  getItem: (key) => window.localStorage.getItem(key),
  setItem: (key, value) => window.localStorage.setItem(key, value),
};
