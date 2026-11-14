import { type ApiClient, ApiError } from "../../../lib/api/client";
import type { AuthGateway } from "../application/ports";
import { type Me, meSchema } from "../domain/me";

/** A 401 of /api/me is an answer (no session), not a failure. */
async function readMe(api: ApiClient): Promise<Me | null> {
  try {
    return meSchema.parse(await api.get("/api/me"));
  } catch (error) {
    if (error instanceof ApiError && error.isUnauthorized) {
      return null;
    }
    throw error;
  }
}

export function createHttpAuthGateway(api: ApiClient): AuthGateway {
  return {
    getMe: () => readMe(api),
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
    logout: async () => {
      await api.send({ method: "POST", path: "/api/auth/logout" });
    },
  };
}
