import type { ApiClient } from "../../../lib/api/client";
import type { AuthGateway } from "../application/ports";
import { meSchema } from "../domain/me";

export function createHttpAuthGateway(api: ApiClient): AuthGateway {
  return {
    getMe: async () => meSchema.parse(await api.get("/api/me")),
    register: async (registration) => {
      await api.send({ method: "POST", path: "/api/auth/register", body: registration });
    },
  };
}
