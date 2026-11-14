import { render, screen } from "@testing-library/react";
import { expect, it } from "vitest";
import { App } from "./App";

it("mounts the app", () => {
  render(<App />);

  expect(screen.getByRole("heading", { name: "Fokko Subs" })).toBeInTheDocument();
});

it("reaches no network in tests", async () => {
  await expect(fetch("/api/health")).rejects.toThrow("fetch is not available in tests");
});
