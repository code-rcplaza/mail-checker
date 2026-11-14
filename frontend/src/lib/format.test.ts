import { describe, expect, it } from "vitest";
import { formatMoney, formatMonth, formatMonths } from "./format";

describe("formatMoney", () => {
  it.each([
    [8_406_000, "CLP", "$84.060"],
    [599_000, "CLP", "$5.990"],
    [1_999, "USD", "US$19,99"],
  ])("writes %d %s cents as %s", (cents, currency, text) => {
    expect(formatMoney({ amount: cents, currency })).toBe(text);
  });
});

describe("formatMonth", () => {
  it("writes a YYYY-MM month in Spanish", () => {
    expect(formatMonth("2026-09")).toBe("septiembre de 2026");
  });

  it("leaves anything else as it came", () => {
    expect(formatMonth("pronto")).toBe("pronto");
  });
});

describe("formatMonths", () => {
  it.each([
    [["2026-10"], "octubre de 2026"],
    [["2026-09", "2026-10"], "septiembre de 2026 y octubre de 2026"],
    [["2026-08", "2026-09", "2026-10"], "agosto de 2026, septiembre de 2026 y octubre de 2026"],
  ])("joins %j", (months, text) => {
    expect(formatMonths(months)).toBe(text);
  });
});
