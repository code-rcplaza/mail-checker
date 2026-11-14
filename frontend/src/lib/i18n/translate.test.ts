import { describe, expect, it } from "vitest";
import { es } from "./es";
import { t, translateServerKey } from "./translate";

describe("t", () => {
  it("returns the text of a key", () => {
    expect(t("common.not_found")).toBe("No encontramos lo que buscas.");
  });

  it("fills the params", () => {
    expect(t("auth.too_many_sends", { limit: 3 })).toBe(
      "Ya enviamos 3 códigos en la última hora. Espera un poco.",
    );
  });

  it("rejects an unknown key at compile time", () => {
    // Never called: the check is that typecheck fails without the directive.
    // @ts-expect-error: not in the catalog
    const unknownKey = () => t("common.nope");
    expect(typeof unknownKey).toBe("function");
  });
});

describe("translateServerKey", () => {
  it("shows the generic error for a key the catalog lacks, never the raw key", () => {
    expect(translateServerKey("billing.brand_new")).toBe(es["common.internal"]);
  });

  it("translates a known key with its params and ignores params it cannot print", () => {
    expect(translateServerKey("auth.name_too_long", { max: 60, extra: { nested: true } })).toBe(
      "El nombre puede tener hasta 60 caracteres.",
    );
  });
});

describe("catalog", () => {
  it("has a message for every common key", () => {
    const common = Object.entries(es).filter(([key]) => key.startsWith("common."));
    expect(common.length).toBeGreaterThan(0);
    for (const [key, text] of common) {
      expect(text.trim(), key).not.toBe("");
    }
  });

  it("gives two keys the same text only on purpose", () => {
    const allowedTwins = new Set<string>();
    const byText = new Map<string, string[]>();
    for (const [key, text] of Object.entries(es)) {
      byText.set(text, [...(byText.get(text) ?? []), key]);
    }
    const twins = [...byText.values()].filter(
      (keys) => keys.length > 1 && !keys.every((key) => allowedTwins.has(key)),
    );
    expect(twins).toEqual([]);
  });
});
