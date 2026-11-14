import { readFileSync } from "node:fs";
import { join } from "node:path";
import { screen, waitFor } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { gatewaysWith } from "../lib/gateways/renderWithGateways";
import { renderApp } from "./renderApp";

const ana = { id: "u1", name: "Ana Pérez", email: "ana@example.com", initials: "AP" };

describe("app router", () => {
  it.each([
    ["with a session", "/resumen", async () => ana],
    ["without a session", "/entrar", async () => null],
    ["when /api/me cannot be read", "/entrar", async () => Promise.reject(new Error("offline"))],
  ])("sends the root %s to %s", async (_case, destination, getMe) => {
    const { router } = renderApp("/", { auth: { ...gatewaysWith({}).auth, getMe } });

    await waitFor(() => expect(router.state.location.pathname).toBe(destination));
  });

  it("shows the 404 screen with a way home for an unknown route", async () => {
    renderApp("/no-existe");

    expect(await screen.findByRole("heading", { name: "No encontramos esta página" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Volver al inicio" })).toHaveAttribute("href", "/");
  });
});

describe("state", () => {
  it("keeps server data in TanStack Query, with no global store library", () => {
    const manifest: { dependencies?: Record<string, string> } = JSON.parse(
      readFileSync(join(import.meta.dirname, "..", "..", "package.json"), "utf8"),
    );
    const stores = [
      "redux",
      "@reduxjs/toolkit",
      "react-redux",
      "zustand",
      "jotai",
      "recoil",
      "mobx",
      "valtio",
    ];

    expect(Object.keys(manifest.dependencies ?? {}).filter((name) => stores.includes(name))).toEqual([]);
  });
});
