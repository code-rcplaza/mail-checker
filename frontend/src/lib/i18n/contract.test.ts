import { readdirSync, readFileSync, statSync } from "node:fs";
import { join } from "node:path";
import { expect, it } from "vitest";
import { es } from "./es";

const backendInternal = join(import.meta.dirname, "../../../../backend/internal");

function goFiles(dir: string): string[] {
  return readdirSync(dir).flatMap((name) => {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) {
      return goFiles(path);
    }
    return name.endsWith(".go") && !name.endsWith("_test.go") ? [path] : [];
  });
}

/** Every `Key… MessageKey = "<feature>.<reason>"` the back declares. */
function backendKeys(): string[] {
  const keys = goFiles(backendInternal).flatMap((file) =>
    [...readFileSync(file, "utf8").matchAll(/MessageKey\s*=\s*"([a-z_]+\.[a-z_]+)"/g)].map(
      (match) => match[1] ?? "",
    ),
  );
  return [...new Set(keys)].sort();
}

it("has a Spanish text for every key of the back", () => {
  const keys = backendKeys();
  // If the way keys are declared changes, finding none must fail loudly.
  expect(keys.length).toBeGreaterThan(50);

  const missing = keys.filter((key) => !Object.hasOwn(es, key));

  expect(missing, `keys without text in es.ts: ${missing.join(", ")}`).toEqual([]);
});
