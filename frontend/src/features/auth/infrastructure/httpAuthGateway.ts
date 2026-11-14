import type { ApiClient } from "../../../lib/api/client";
import type { AuthGateway } from "../application/ports";
import { meSchema } from "../domain/me";

export function createHttpAuthGateway(api: ApiClient): AuthGateway {
  return {
    getMe: async () => meSchema.parse(await api.get("/api/me")),
  };
}
