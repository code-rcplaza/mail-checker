/** The ring of the handoff: 240 × 240 box, radius 100, stroke 8, 14 of gap between arcs. */
export const ringRadius = 100;
export const ringCircumference = 2 * Math.PI * ringRadius;
export const ringGap = 14;

export interface RingArc {
  /** Length of the drawn arc along the circle. */
  readonly length: number;
  /** Where it starts, as an SVG dash offset (negative: clockwise from the top). */
  readonly offset: number;
}

/**
 * Each value gets its share of the circle minus the gap; a single value
 * closes the circle, and nothing to show (all zero) draws nothing.
 */
export function ringArcs(values: ReadonlyArray<number>): RingArc[] {
  const positive = values.filter((value) => value > 0);
  const total = positive.reduce((sum, value) => sum + value, 0);
  if (total === 0) {
    return [];
  }
  if (positive.length === 1) {
    return [{ length: ringCircumference, offset: 0 }];
  }
  let start = 0;
  return positive.map((value) => {
    const share = (value / total) * ringCircumference;
    const arc = { length: Math.max(share - ringGap, 0), offset: -start };
    start += share;
    return arc;
  });
}
