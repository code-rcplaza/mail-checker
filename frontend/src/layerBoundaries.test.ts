import { execFileSync } from "node:child_process";
import { mkdirSync, rmSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { afterAll, describe, expect, it } from "vitest";

// Each case writes a probe file in a layer, lints it with the project's Biome
// config and checks whether noRestrictedImports fires. Probes are removed after.
const projectDir = join(import.meta.dirname, "..");
const probeFeature = "src/features/zzlayerprobe";
const probeComponent = "src/components/ui/atoms/ZzLayerProbe.tsx";
const probeLib = "src/lib/zzLayerProbe.ts";

function lintFires(relativePath: string, source: string): boolean {
  const path = join(projectDir, relativePath);
  mkdirSync(dirname(path), { recursive: true });
  writeFileSync(path, source);
  return lintReports(relativePath);
}

/** Lints a file as it is, without writing it. Biome prints its diagnostics on stderr. */
function lintReports(relativePath: string): boolean {
  try {
    execFileSync(join(projectDir, "node_modules/.bin/biome"), ["lint", relativePath], {
      cwd: projectDir,
      stdio: "pipe",
    });
    return false;
  } catch (failure) {
    const output = failure instanceof Error && "stderr" in failure ? String(failure.stderr) : "";
    return output.includes("noRestrictedImports");
  }
}

afterAll(() => {
  rmSync(join(projectDir, probeFeature), { recursive: true, force: true });
  rmSync(join(projectDir, probeComponent), { force: true });
  rmSync(join(projectDir, probeLib), { force: true });
});

const importOf = (specifier: string) =>
  `import * as probe from "${specifier}";\nexport const used = probe;\n`;

describe("layer boundaries", () => {
  it.each([
    ["react", `${probeFeature}/domain/a.ts`, "react"],
    ["TanStack", `${probeFeature}/domain/b.ts`, "@tanstack/react-query"],
    ["application", `${probeFeature}/domain/c.ts`, "../application/ports"],
  ])("domain cannot import %s", (_what, file, specifier) => {
    expect(lintFires(file, importOf(specifier))).toBe(true);
  });

  it.each([
    [
      "TanStack from a feature component",
      `${probeFeature}/presentation/components/A.tsx`,
      "@tanstack/react-query",
    ],
    [
      "application from a feature component",
      `${probeFeature}/presentation/components/B.tsx`,
      "../../application/ports",
    ],
    [
      "infrastructure from a feature component",
      `${probeFeature}/presentation/components/C.tsx`,
      "../../infrastructure/x",
    ],
    ["TanStack from components/ui", probeComponent, "@tanstack/react-router"],
    ["application from components/ui", probeComponent, "../../../features/auth/application/ports"],
  ])("blocks %s", (_what, file, specifier) => {
    expect(lintFires(file, importOf(specifier))).toBe(true);
  });

  it("lets a presentational piece import react and lib/i18n", () => {
    const source = `import { useId } from "react";\nimport { t } from "../../../lib/i18n/translate";\nexport const used = [useId, t];\n`;

    expect(lintFires(probeComponent, source)).toBe(false);
  });

  it("lets only main.tsx and infrastructure import infrastructure", () => {
    expect(lintFires(probeLib, importOf("../features/auth/infrastructure/httpAuthGateway"))).toBe(true);
    expect(
      lintFires(
        `${probeFeature}/infrastructure/other.ts`,
        importOf("../../auth/infrastructure/httpAuthGateway"),
      ),
    ).toBe(false);
    expect(lintReports("src/main.tsx")).toBe(false);
  });
});
