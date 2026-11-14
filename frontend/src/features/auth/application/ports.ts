import type { Me } from "../domain/me";
import type { Registration } from "../domain/registration";

/** What the auth feature needs from the outside; tests give a double. */
export interface AuthGateway {
  getMe(): Promise<Me>;
  /** Sends the code to the email; the answer is the same whether the account existed or not. */
  register(registration: Registration): Promise<void>;
}
