import { readdirSync, readFileSync, statSync } from "node:fs";
import { join, relative } from "node:path";
import { expect, it } from "vitest";

const srcDir = join(import.meta.dirname, "..");

/** The only place where raw colors live. */
const tokensFile = join(srcDir, "styles", "tokens.css");

function sourceFiles(dir: string): string[] {
  return readdirSync(dir).flatMap((name) => {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) {
      return sourceFiles(path);
    }
    return /\.(tsx?|css)$/.test(name) && !/\.test\.tsx?$/.test(name) ? [path] : [];
  });
}

const forbidden = [
  { what: "a raw token", pattern: /--raw-/ },
  { what: "a hex color", pattern: /#[0-9a-fA-F]{3,8}\b/ },
  { what: "an rgb/hsl/oklch color", pattern: /\b(?:rgba?|hsla?|oklch)\(/ },
  { what: "a raw token class", pattern: /\b(?:bg|text|border|fill|stroke|ring)-raw-/ },
];

it("uses only semantic colors outside tokens.css", () => {
  const offenders: string[] = [];
  for (const file of sourceFiles(srcDir).filter((path) => path !== tokensFile)) {
    const content = readFileSync(file, "utf8");
    for (const { what, pattern } of forbidden) {
      if (pattern.test(content)) {
        offenders.push(`${relative(srcDir, file)} uses ${what}`);
      }
    }
  }

  expect(offenders).toEqual([]);
});
