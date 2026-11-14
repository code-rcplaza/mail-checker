import { describe, expect, it } from "vitest";
import { billingDateToUnix, fileProblemOf, maxUploadBytes, uploadFormOf } from "./upload";

describe("fileProblemOf", () => {
  it.each([
    ["a PDF", { name: "cartola.pdf", type: "application/pdf", size: 1000 }, undefined],
    ["a PDF with no type", { name: "CARTOLA.PDF", type: "", size: 1000 }, undefined],
    ["an image", { name: "foto.png", type: "image/png", size: 1000 }, "statements.file_kind_unknown"],
    ["a CSV", { name: "datos.csv", type: "text/csv", size: 1000 }, "statements.file_kind_unknown"],
    ["an empty PDF", { name: "cartola.pdf", type: "application/pdf", size: 0 }, "statements.file_empty"],
    ["exactly 10 MB", { name: "cartola.pdf", type: "application/pdf", size: maxUploadBytes }, undefined],
    [
      "over 10 MB",
      { name: "cartola.pdf", type: "application/pdf", size: maxUploadBytes + 1 },
      "statements.file_too_large",
    ],
  ])("judges %s", (_case, file, problem) => {
    expect(fileProblemOf(file)).toBe(problem);
  });
});

describe("billingDateToUnix", () => {
  it("reads the day of a date input as unix seconds at noon UTC, so no time zone moves the day", () => {
    expect(billingDateToUnix("2026-10-15")).toBe(Date.UTC(2026, 9, 15, 12) / 1000);
  });

  it.each(["", "15-10-2026", "2026-02-30", "mañana"])("reads nothing from %j", (text) => {
    expect(billingDateToUnix(text)).toBeUndefined();
  });
});

describe("uploadFormOf", () => {
  const file = new File(["%PDF"], "cartola.pdf", { type: "application/pdf" });

  it("sends only the file when bank and date are not given", () => {
    const form = uploadFormOf({ file });

    expect([...form.keys()]).toEqual(["file"]);
  });

  it("adds the bank and the billing date that are present", () => {
    const form = uploadFormOf({ file, bank: "bci", billingDate: 1_760_529_600 });

    expect(form.get("file")).toBe(file);
    expect(form.get("bank")).toBe("bci");
    expect(form.get("billingDate")).toBe("1760529600");
  });
});
