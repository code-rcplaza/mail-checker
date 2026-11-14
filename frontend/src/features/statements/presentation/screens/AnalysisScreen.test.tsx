import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { renderApp } from "../../../../app/renderApp";
import { gatewaysWith } from "../../../../lib/gateways/renderWithGateways";
import type { StatementsGateway } from "../../application/ports";
import { uploadedStatementKey } from "../../application/useUploadStatement";

const ana = { id: "u1", name: "Ana Pérez", email: "ana@example.com", initials: "AP" };
const uploaded = {
  id: "s1",
  filename: "cartola-octubre.pdf",
  months: ["2026-09", "2026-10"],
  charges: 12,
  badLines: 0,
};

function renderAnalysis(
  found: { sure: number; possible: number },
  options: { withUpload: boolean } = { withUpload: true },
) {
  const statements: StatementsGateway = {
    ...gatewaysWith({}).statements,
    readFound: vi.fn(async () => ({
      ...found,
      monthly: { amount: 8_406_000, currency: "CLP", isApproximate: false },
    })),
  };
  const auth = { ...gatewaysWith({}).auth, getMe: async () => ana };
  const rendered = renderApp("/analisis/s1", { auth, statements });
  if (options.withUpload) {
    rendered.queryClient.setQueryData(uploadedStatementKey("s1"), uploaded);
  }
  return rendered;
}

describe("Analysis screen", () => {
  it("shows the file and the months it read", async () => {
    renderAnalysis({ sure: 3, possible: 2 });

    expect(await screen.findByText("Leímos cartola-octubre.pdf")).toBeInTheDocument();
    expect(screen.getByText("Meses: septiembre de 2026 y octubre de 2026")).toBeInTheDocument();
  });

  it("shows how many are sure and how many need a look", async () => {
    renderAnalysis({ sure: 3, possible: 2 });

    expect(await screen.findByText("3 seguras")).toBeInTheDocument();
    expect(screen.getByText("2 por confirmar")).toBeInTheDocument();
    expect(screen.getByText("$84.060 al mes.")).toBeInTheDocument();
  });

  it("goes to Revisión when some need a look", async () => {
    const { router } = renderAnalysis({ sure: 3, possible: 2 });

    await userEvent.click(await screen.findByRole("link", { name: "Revisar las 2 dudas" }));

    await waitFor(() => expect(router.state.location.pathname).toBe("/revision"));
  });

  it("goes to the Resumen when none needs a look", async () => {
    const { router } = renderAnalysis({ sure: 3, possible: 0 });

    await userEvent.click(await screen.findByRole("link", { name: "Ver mi resumen" }));

    await waitFor(() => expect(router.state.location.pathname).toBe("/resumen"));
  });

  it("says only what is really kept", async () => {
    renderAnalysis({ sure: 1, possible: 0 });

    const kept = await screen.findByRole("list", { name: "Guardamos" });
    expect(kept).toHaveTextContent("Comercio");
    expect(kept).toHaveTextContent("Últimos 4 dígitos de la tarjeta");
    expect(screen.getByRole("list", { name: "No guardamos" })).toHaveTextContent("El PDF ni su texto");
  });

  it("offers another statement when nothing was found", async () => {
    const { router } = renderAnalysis({ sure: 0, possible: 0 });

    expect(await screen.findByRole("heading", { name: "No encontramos suscripciones" })).toBeInTheDocument();
    await userEvent.click(screen.getByRole("link", { name: "Subir otra" }));

    await waitFor(() => expect(router.state.location.pathname).toBe("/fuente"));
  });

  it("still shows the counts after a reload, without the file", async () => {
    renderAnalysis({ sure: 3, possible: 2 }, { withUpload: false });

    expect(await screen.findByText("3 seguras")).toBeInTheDocument();
    expect(screen.queryByText(/Leímos/)).not.toBeInTheDocument();
  });
});
