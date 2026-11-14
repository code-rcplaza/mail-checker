import { act, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { AnimatedNumber, countDuration } from "./AnimatedNumber";

// jsdom has no matchMedia: each test sets the preference it needs and the original comes back after.
const originalMatchMedia = Object.getOwnPropertyDescriptor(window, "matchMedia");

function preferReducedMotion(reduce: boolean) {
  Object.defineProperty(window, "matchMedia", {
    configurable: true,
    value: (query: string) => ({ matches: reduce && query.includes("reduce"), media: query }),
  });
}

function restoreMatchMedia() {
  if (originalMatchMedia === undefined) {
    Reflect.deleteProperty(window, "matchMedia");
    return;
  }
  Object.defineProperty(window, "matchMedia", originalMatchMedia);
}

function advance(ms: number) {
  act(() => vi.advanceTimersByTime(ms));
}

const shownNumber = () => Number(screen.getByText(/^\d+$/).textContent);

describe("AnimatedNumber", () => {
  beforeEach(() => {
    vi.useFakeTimers({
      toFake: ["requestAnimationFrame", "cancelAnimationFrame", "performance", "setTimeout"],
    });
    preferReducedMotion(false);
  });
  afterEach(() => {
    vi.useRealTimers();
    restoreMatchMedia();
  });

  it("counts from 0 to the value in 0.9 s and ends exactly on it", () => {
    render(<AnimatedNumber value={12990} />);
    expect(shownNumber()).toBe(0);

    advance(countDuration / 2);
    expect(shownNumber()).toBeGreaterThan(0);
    expect(shownNumber()).toBeLessThan(12990);
    advance(countDuration);

    expect(shownNumber()).toBe(12990);
  });

  it("shows the final value at once with reduced motion", () => {
    preferReducedMotion(true);

    render(<AnimatedNumber value={12990} />);

    expect(shownNumber()).toBe(12990);
  });

  it("animates from the shown value to a new one", () => {
    const { rerender } = render(<AnimatedNumber value={1000} />);
    advance(countDuration * 2);

    rerender(<AnimatedNumber value={2000} />);
    advance(countDuration / 2);

    expect(shownNumber()).toBeGreaterThan(1000);
    expect(shownNumber()).toBeLessThan(2000);
    advance(countDuration);
    expect(shownNumber()).toBe(2000);
  });

  it("formats every frame", () => {
    preferReducedMotion(true);

    render(<AnimatedNumber value={7990} format={(value) => `$${value}`} />);

    expect(screen.getByText("$7990")).toBeInTheDocument();
  });
});
