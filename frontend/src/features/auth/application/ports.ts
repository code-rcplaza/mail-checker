import type { Me } from "../domain/me";

/** What the auth feature needs from the outside; tests give a double. */
export interface AuthGateway {
  getMe(): Promise<Me>;
}
