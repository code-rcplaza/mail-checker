import type { ApiClient } from "../../../lib/api/client";
import type { StatementsGateway } from "../application/ports";
import { banksSchema, foundSchema, uploadFormOf, uploadResultSchema } from "../domain/upload";

export function createHttpStatementsGateway(api: ApiClient): StatementsGateway {
  return {
    listBanks: async () => banksSchema.parse(await api.get("/api/banks")),
    upload: async (input) =>
      uploadResultSchema.parse(
        await api.send({ method: "POST", path: "/api/statements", body: uploadFormOf(input) }),
      ),
    readFound: async () => foundSchema.parse(await api.get("/api/summary")),
  };
}
