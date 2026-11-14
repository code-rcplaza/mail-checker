import type { ReactNode } from "react";
import { classNames } from "../../../lib/classNames";
import { t } from "../../../lib/i18n/translate";
import { ringArcs, ringCircumference, ringRadius } from "./ringGeometry";

export type RingSize = 224 | 176 | 144 | 120 | 104;
export type SeriesIndex = 1 | 2 | 3 | 4 | 5 | 6;

export interface RingSegment {
  readonly value: number;
  readonly series: SeriesIndex;
}

export type SegmentedRingProps =
  | {
      readonly variant: "normal";
      readonly segments: ReadonlyArray<RingSegment>;
      readonly label: string;
      readonly size: RingSize;
      readonly children?: ReactNode;
    }
  | {
      readonly variant: "empty";
      readonly label: string;
      readonly size: RingSize;
      readonly children?: ReactNode;
    }
  | { readonly variant: "loading"; readonly size: RingSize; readonly children?: ReactNode };

const strokeBySeries: Readonly<Record<SeriesIndex, string>> = {
  1: "stroke-series-1",
  2: "stroke-series-2",
  3: "stroke-series-3",
  4: "stroke-series-4",
  5: "stroke-series-5",
  6: "stroke-series-6",
};

/** Tailwind sizes are 4 px steps: 224 px is size-56. */
const sizeClasses: Readonly<Record<RingSize, string>> = {
  224: "size-56",
  176: "size-44",
  144: "size-36",
  120: "size-30",
  104: "size-26",
};

const center = 120;
const rotateFromTop = `rotate(-90 ${center} ${center})`;

function SpendArcs({ segments }: { readonly segments: ReadonlyArray<RingSegment> }) {
  const shown = segments.filter((segment) => segment.value > 0);
  const arcs = ringArcs(shown.map((segment) => segment.value));
  return arcs.map((arc, index) => (
    <circle
      // biome-ignore lint/suspicious/noArrayIndexKey: arcs follow the order of the segments
      key={index}
      data-arc=""
      cx={center}
      cy={center}
      r={ringRadius}
      fill="none"
      strokeWidth={8}
      strokeLinecap="round"
      strokeDasharray={`${arc.length.toFixed(1)} ${ringCircumference.toFixed(1)}`}
      strokeDashoffset={arc.offset.toFixed(1)}
      transform={rotateFromTop}
      className={strokeBySeries[shown[index]?.series ?? 6]}
    />
  ));
}

function ringContent(props: SegmentedRingProps): ReactNode {
  switch (props.variant) {
    case "normal":
      return <SpendArcs segments={props.segments} />;
    case "empty":
      return (
        <circle
          data-empty=""
          cx={center}
          cy={center}
          r={ringRadius}
          fill="none"
          strokeWidth={4}
          strokeDasharray="2 10"
          strokeLinecap="round"
          className="stroke-border-strong"
        />
      );
    case "loading":
      return (
        <g className="origin-center motion-safe:animate-spin" data-loading="">
          <circle
            cx={center}
            cy={center}
            r={ringRadius}
            fill="none"
            strokeWidth={8}
            className="stroke-border"
          />
          <circle
            cx={center}
            cy={center}
            r={ringRadius}
            fill="none"
            strokeWidth={8}
            strokeLinecap="round"
            strokeDasharray={`${(ringCircumference / 4).toFixed(1)} ${ringCircumference.toFixed(1)}`}
            transform={rotateFromTop}
            className="stroke-primary"
          />
        </g>
      );
  }
}

/** The visual signature: the spend split by service, with a thin outer ring. */
export function SegmentedRing(props: SegmentedRingProps) {
  const isLoading = props.variant === "loading";
  return (
    <div className={classNames("relative inline-flex items-center justify-center", sizeClasses[props.size])}>
      <svg
        className="size-full"
        viewBox="0 0 240 240"
        role={isLoading ? "status" : "img"}
        aria-label={isLoading ? t("ui.loading") : props.label}
      >
        <circle
          cx={center}
          cy={center}
          r={115}
          fill="none"
          strokeWidth={1.5}
          className={classNames("stroke-hero-line", "opacity-70")}
        />
        {ringContent(props)}
      </svg>
      {props.children !== undefined && (
        <div className="absolute inset-0 flex flex-col items-center justify-center">{props.children}</div>
      )}
    </div>
  );
}
