import { describe, expect, it, vi } from "vitest";
import { ApiError, createApiClient } from "../../../lib/api/client";
import { createHttpStatementsGateway } from "./httpStatementsGateway";

function fetchAnswering(status: number, body: unknown) {
  return vi.fn(async (_input: string, _init: RequestInit) => new Response(JSON.stringify(body), { status }));
}

const file = new File(["%PDF"], "cartola.pdf", { type: "application/pdf" });
const uploaded = {
  id: "s1",
  filename: "cartola.pdf",
  months: ["2026-09", "2026-10"],
  charges: 42,
  badLines: 0,
};

describe("http statements gateway", () => {
  it("sends the file and the bank as multipart to /api/statements and reads the result", async () => {
    const fetchFunction = fetchAnswering(201, uploaded);

    const result = await createHttpStatementsGateway(createApiClient(fetchFunction)).upload({
      file,
      bank: "bci",
    });

    expect(result).toEqual(uploaded);
    const [path, init] = fetchFunction.mock.calls[0] ?? [];
    expect(path).toBe("/api/statements");
    expect(init?.method).toBe("POST");
    expect(init?.body).toBeInstanceOf(FormData);
  });

  it("throws ApiError with the key on a 422", async () => {
    const gateway = createHttpStatementsGateway(
      createApiClient(fetchAnswering(422, { error: { key: "statements.no_charges" } })),
    );

    const error = await gateway.upload({ file }).catch((caught: unknown) => caught);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({ status: 422, key: "statements.no_charges" });
  });

  it("lists the banks of /api/banks", async () => {
    const gateway = createHttpStatementsGateway(createApiClient(fetchAnswering(200, [{ id: "bci" }])));

    await expect(gateway.listBanks()).resolves.toEqual([{ id: "bci" }]);
  });

  it("counts the sure and the possible subscriptions of /api/summary", async () => {
    const summary = {
      subscriptions: [{ lineKey: "a" }, { lineKey: "b" }, { lineKey: "c" }],
      possible: [{ lineKey: "d" }, { lineKey: "e" }],
      total: { amount: 8_406_000, currency: "CLP", isApproximate: false, isIncomplete: false },
    };
    const fetchFunction = fetchAnswering(200, summary);

    await expect(createHttpStatementsGateway(createApiClient(fetchFunction)).readFound()).resolves.toEqual({
      sure: 3,
      possible: 2,
      monthly: { amount: 8_406_000, currency: "CLP", isApproximate: false },
    });
    expect(fetchFunction.mock.calls[0]?.[0]).toBe("/api/summary");
  });
});
