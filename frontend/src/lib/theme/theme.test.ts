import { describe, expect, it } from "vitest";
import {
  applyThemePreference,
  type PreferenceStorage,
  readThemePreference,
  storeThemePreference,
} from "./theme";

function memoryStorage(): PreferenceStorage {
  const values = new Map<string, string>();
  return {
    getItem: (key) => values.get(key) ?? null,
    setItem: (key, value) => {
      values.set(key, value);
    },
  };
}

const blockedStorage: PreferenceStorage = {
  getItem: () => {
    throw new DOMException("blocked", "SecurityError");
  },
  setItem: () => {
    throw new DOMException("full", "QuotaExceededError");
  },
};

describe("theme preference", () => {
  it("applies dark and light over the system, and leaves system to the media query", () => {
    const root = document.createElement("html");

    applyThemePreference(root, "dark");
    expect(root.dataset.theme).toBe("dark");
    applyThemePreference(root, "light");
    expect(root.dataset.theme).toBe("light");
    applyThemePreference(root, "system");
    expect(root.hasAttribute("data-theme")).toBe(false);
  });

  it("keeps the chosen preference across reloads", () => {
    const storage = memoryStorage();

    expect(storeThemePreference(storage, "dark")).toBe(true);

    expect(readThemePreference(storage)).toBe("dark");
  });

  it("falls back to system without a stored or with an unknown value", () => {
    const storage = memoryStorage();
    expect(readThemePreference(storage)).toBe("system");

    storage.setItem("fokko.theme", "sepia");
    expect(readThemePreference(storage)).toBe("system");
  });

  it("works with blocked or full storage", () => {
    expect(readThemePreference(blockedStorage)).toBe("system");
    expect(storeThemePreference(blockedStorage, "dark")).toBe(false);
  });
});
