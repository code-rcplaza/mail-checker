import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { ringArcs, ringCircumference, ringGap } from "./ringGeometry";
import { type RingSize, SegmentedRing } from "./SegmentedRing";

describe("ringArcs", () => {
  it("gives each value its share of the circle minus the gap", () => {
    const arcs = ringArcs([50, 30, 20]);

    expect(arcs.map((arc) => arc.length)).toEqual(
      [0.5, 0.3, 0.2]
        .map((share) => share * ringCircumference - ringGap)
        .map((length) => expect.closeTo(length, 6)),
    );
    expect(arcs.map((arc) => arc.offset)).toEqual(
      [0, -0.5, -0.8].map((start) => expect.closeTo(start * ringCircumference, 6)),
    );
  });

  it("closes the circle with a single value", () => {
    expect(ringArcs([12990])).toEqual([{ length: ringCircumference, offset: 0 }]);
  });

  it("draws nothing for a zero total, without dividing by zero", () => {
    expect(ringArcs([])).toEqual([]);
    expect(ringArcs([0, 0])).toEqual([]);
  });
});

describe("SegmentedRing", () => {
  it("draws one arc per segment with its series color", () => {
    const { container } = render(
      <SegmentedRing
        variant="normal"
        size={224}
        label="Tu gasto del mes por servicio"
        segments={[
          { value: 50, series: 1 },
          { value: 30, series: 3 },
        ]}
      />,
    );

    expect(screen.getByRole("img", { name: "Tu gasto del mes por servicio" })).toBeInTheDocument();
    const arcs = container.querySelectorAll("[data-arc]");
    expect([...arcs].map((arc) => arc.getAttribute("class"))).toEqual(["stroke-series-1", "stroke-series-3"]);
  });

  it("draws a dotted ring when empty", () => {
    const { container } = render(<SegmentedRing variant="empty" size={144} label="Sin gasto aún" />);

    expect(container.querySelector("[data-empty]")).toHaveAttribute("stroke-dasharray", "2 10");
  });

  it("announces loading and spins only when motion is allowed", () => {
    const { container } = render(<SegmentedRing variant="loading" size={120} />);

    expect(screen.getByRole("status", { name: "Cargando" })).toBeInTheDocument();
    const spinner = container.querySelector("[data-loading]");
    expect(spinner).toHaveClass("motion-safe:animate-spin");
    expect(spinner).not.toHaveClass("animate-spin");
  });

  it.each([
    [224, "size-56"],
    [176, "size-44"],
    [144, "size-36"],
    [120, "size-30"],
    [104, "size-26"],
  ] as const)("measures %i px", (size: RingSize, sizeClass) => {
    const { container } = render(<SegmentedRing variant="empty" size={size} label="Anillo" />);

    expect(container.firstElementChild).toHaveClass(sizeClass);
  });
});
