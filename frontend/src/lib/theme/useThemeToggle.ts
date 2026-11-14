import { useState } from "react";
import {
  applyThemePreference,
  browserStorage,
  readThemePreference,
  storeThemePreference,
  type ThemePreference,
} from "./theme";

/** "system" follows the browser; light and dark are what they say. */
export function isDarkTheme(preference: ThemePreference, systemIsDark: boolean): boolean {
  return preference === "dark" || (preference === "system" && systemIsDark);
}

function systemIsDark(): boolean {
  return typeof window.matchMedia === "function" && window.matchMedia("(prefers-color-scheme: dark)").matches;
}

/** The moon/sun button of the access screens: it switches to the other theme and keeps it. */
export function useThemeToggle() {
  const [preference, setPreference] = useState(() => readThemePreference(browserStorage));
  const isDark = isDarkTheme(preference, systemIsDark());

  function onToggle() {
    const next: ThemePreference = isDark ? "light" : "dark";
    storeThemePreference(browserStorage, next);
    applyThemePreference(document.documentElement, next);
    setPreference(next);
  }

  return { isDark, onToggle };
}
