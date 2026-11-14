import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

// The Go server embeds the build, so it lands inside the backend and the app
// ships as one binary. In development Vite serves the front and sends /api to Go.
const backendURL = "http://127.0.0.1:8080";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    strictPort: true,
    proxy: { "/api": { target: backendURL, changeOrigin: false } },
  },
  build: {
    outDir: "../backend/internal/shared/httpserver/web/dist",
    emptyOutDir: true,
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test/setup.ts"],
    css: false,
  },
});
