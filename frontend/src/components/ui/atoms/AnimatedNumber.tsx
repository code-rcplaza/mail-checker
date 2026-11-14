import { useEffect, useRef, useState } from "react";

export const countDuration = 900;

function prefersReducedMotion(): boolean {
  return (
    typeof window.matchMedia === "function" && window.matchMedia("(prefers-reduced-motion: reduce)").matches
  );
}

function easeOutCubic(progress: number): number {
  return 1 - (1 - progress) ** 3;
}

export interface AnimatedNumberProps {
  readonly value: number;
  /** How to show each frame; whole numbers by default. */
  readonly format?: (value: number) => string;
}

const roundToString = (value: number) => String(Math.round(value));

/** Counts to the value in 0.9 s (from 0 on mount, from the shown one after); still with reduced motion. */
export function AnimatedNumber({ value, format = roundToString }: AnimatedNumberProps) {
  const [shown, setShown] = useState(() => (prefersReducedMotion() ? value : 0));
  const shownRef = useRef(shown);

  useEffect(() => {
    const from = shownRef.current;
    if (prefersReducedMotion() || from === value) {
      shownRef.current = value;
      setShown(value);
      return;
    }
    let frame = 0;
    const startedAt = performance.now();
    const step = (now: number) => {
      const progress = Math.min((now - startedAt) / countDuration, 1);
      const next = progress === 1 ? value : from + (value - from) * easeOutCubic(progress);
      shownRef.current = next;
      setShown(next);
      if (progress < 1) {
        frame = requestAnimationFrame(step);
      }
    };
    frame = requestAnimationFrame(step);
    return () => cancelAnimationFrame(frame);
  }, [value]);

  return <span className="tabular-nums">{format(shown)}</span>;
}
