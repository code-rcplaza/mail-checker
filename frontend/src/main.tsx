import "@fontsource/be-vietnam-pro/400.css";
import "@fontsource/be-vietnam-pro/500.css";
import "@fontsource/be-vietnam-pro/600.css";
import "@fontsource/prompt/600.css";
import "./styles/index.css";

import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { AppProviders } from "./app/AppProviders";
import { createAppRouter } from "./app/router";
import { createHttpAuthGateway } from "./features/auth/infrastructure/httpAuthGateway";
import { createApiClient } from "./lib/api/client";
import type { Gateways } from "./lib/gateways/gateways";
import { createQueryClient } from "./lib/query/queryClient";
import { applyThemePreference, browserStorage, readThemePreference } from "./lib/theme/theme";

const root = document.getElementById("root");
if (root === null) {
  throw new Error("index.html has no #root element");
}

// Before the first paint, so a dark preference never flashes light.
applyThemePreference(document.documentElement, readThemePreference(browserStorage));

// The composition root: the only place that builds the real gateways.
const api = createApiClient((input, init) => window.fetch(input, init));
const gateways: Gateways = { auth: createHttpAuthGateway(api) };
const queryClient = createQueryClient();
const router = createAppRouter({ queryClient, gateways });

createRoot(root).render(
  <StrictMode>
    <AppProviders queryClient={queryClient} gateways={gateways} router={router} />
  </StrictMode>,
);
