import "@fontsource/be-vietnam-pro/400.css";
import "@fontsource/be-vietnam-pro/500.css";
import "@fontsource/be-vietnam-pro/600.css";
import "@fontsource/prompt/600.css";
import "./styles/index.css";

import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./app/App";
import { applyThemePreference, browserStorage, readThemePreference } from "./lib/theme/theme";

const root = document.getElementById("root");
if (root === null) {
  throw new Error("index.html has no #root element");
}

// Before the first paint, so a dark preference never flashes light.
applyThemePreference(document.documentElement, readThemePreference(browserStorage));

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
