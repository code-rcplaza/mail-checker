import { readFileSync } from "node:fs";
import { join } from "node:path";
import { compile } from "tailwindcss";
import { describe, expect, it } from "vitest";

const stylesDir = join(import.meta.dirname);
const tokensCSS = readFileSync(join(stylesDir, "tokens.css"), "utf8");

/** The color keys of the handoff tokens.json; each one must have a semantic token. */
const handoffColorKeys = [
  "bg",
  "field",
  "navBg",
  "navAct",
  "avBg",
  "skel",
  "line",
  "line2",
  "text",
  "t2",
  "muted",
  "primary",
  "ptext",
  "accent",
  "avT",
  "outline",
  "outlineT",
  "ok",
  "hero",
  "heroT",
  "heroT2",
  "heroLine",
  "heroChip",
  "zT",
  "alBg",
  "alBd",
  "alT",
  "alT2",
  "attn",
  "attnSolid",
  "attnSolidT",
  "attnText",
  "attnBig",
  "destr",
  "destrT",
  "errBg",
  "zone",
];

/** The service colors of the handoff, as data series. */
const handoffSeriesKeys = ["Adobe", "ChatGPT", "Netflix", "iCloud", "Disney\\+", "Otros"];

function blockOf(selector: string): string {
  const start = tokensCSS.indexOf(`${selector} {`);
  if (start < 0) {
    throw new Error(`tokens.css has no ${selector} block`);
  }
  return tokensCSS.slice(start, tokensCSS.indexOf("}", start));
}

function semanticNamesIn(block: string): string[] {
  return [...block.matchAll(/--color-([a-z0-9-]+):\s*var\(--raw-/g)].map((match) => match[1] ?? "");
}

describe("design tokens", () => {
  const light = semanticNamesIn(blockOf("@theme"));
  const dark = semanticNamesIn(blockOf(':root[data-theme="dark"]'));

  it("maps every handoff color to a semantic token", () => {
    for (const key of handoffColorKeys) {
      expect(tokensCSS, `handoff key ${key}`).toMatch(new RegExp(`\\*\\s+${key}\\s+→ [a-z-]+\\n`));
    }
    for (const key of handoffSeriesKeys) {
      expect(tokensCSS, `service color ${key}`).toMatch(new RegExp(`\\*\\s+${key}\\s+→ series-\\d\\n`));
    }
    // One per handoff color, one per series, and the ink on series.
    expect(light.length).toBe(handoffColorKeys.length + handoffSeriesKeys.length + 1);
  });

  it("gives every semantic token a dark value, also for the system preference", () => {
    expect(dark).toEqual(light);
    expect(semanticNamesIn(blockOf(':root:not([data-theme="light"])'))).toEqual(light);
  });

  it("points semantic tokens only to raw steps that exist", () => {
    const defined = new Set([...tokensCSS.matchAll(/(--raw-[a-z0-9-]+):\s*#/g)].map((match) => match[1]));
    for (const used of tokensCSS.matchAll(/var\((--raw-[a-z0-9-]+)\)/g)) {
      expect(defined.has(used[1]), `${used[1]} is not defined`).toBe(true);
    }
  });

  it("makes classes from semantic tokens only", async () => {
    const tailwindDir = join(stylesDir, "../../node_modules/tailwindcss");
    const compiler = await compile(`@import "tailwindcss";\n${tokensCSS}`, {
      base: stylesDir,
      loadStylesheet: async (id: string, base: string) => {
        const path = id === "tailwindcss" ? join(tailwindDir, "index.css") : join(base, id);
        return { path, base: stylesDir, content: readFileSync(path, "utf8") };
      },
    });

    const css = compiler.build([
      "bg-page",
      "text-ink",
      "bg-red-500",
      "text-raw-purple-light-9",
      "bg-raw-onyx-dark-1",
    ]);

    expect(css).toContain("background-color: var(--color-page)");
    expect(css).toContain("color: var(--color-ink)");
    expect(css).not.toContain(".bg-red-500");
    expect(css).not.toContain("raw-purple-light-9 {");
    expect(css).not.toContain(".bg-raw-onyx-dark-1");
  });
});
