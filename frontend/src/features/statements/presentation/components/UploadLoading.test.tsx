import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { UploadLoadingContent, UploadLoadingZone } from "./UploadLoading";

describe("UploadLoading", () => {
  it("shows the loading ring and step 2 of 3", () => {
    const { container } = render(<UploadLoadingZone />);

    expect(screen.getByText("Paso 2 de 3 · Análisis")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Leyendo tus cobros" })).toBeInTheDocument();
    expect(container.querySelector("[data-loading]")).not.toBeNull();
  });

  it("announces the wait and draws placeholder rows", () => {
    const { container } = render(<UploadLoadingContent />);

    expect(screen.getByRole("status")).toHaveTextContent("Leyendo tu cartola");
    expect(
      container.querySelectorAll("[aria-hidden='true'].motion-safe\\:animate-pulse").length,
    ).toBeGreaterThan(0);
  });
});
