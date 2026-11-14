import { readFileSync } from "node:fs";
import { join } from "node:path";
import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { renderApp } from "./renderApp";

describe("app router", () => {
  it("shows the home route", async () => {
    renderApp("/");

    expect(await screen.findByRole("heading", { name: "Fokko Subs" })).toBeInTheDocument();
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
