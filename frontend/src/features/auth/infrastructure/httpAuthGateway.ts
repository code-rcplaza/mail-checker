import type { ApiClient } from "../../../lib/api/client";
import type { AuthGateway } from "../application/ports";
import { meSchema } from "../domain/me";

export function createHttpAuthGateway(api: ApiClient): AuthGateway {
  return {
    getMe: async () => meSchema.parse(await api.get("/api/me")),
    register: async (registration) => {
      await api.send({ method: "POST", path: "/api/auth/register", body: registration });
    },
    login: async (email) => {
      await api.send({ method: "POST", path: "/api/auth/login", body: { email } });
    },
    verifyCode: async (entry) => {
      await api.send({ method: "POST", path: "/api/auth/code/verify", body: entry });
    },
    resendCode: async (email) => {
      await api.send({ method: "POST", path: "/api/auth/code/resend", body: { email } });
    },
  };
}
