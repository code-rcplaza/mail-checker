import type { CodeEntry } from "../domain/code";
import type { Me } from "../domain/me";
import type { Registration } from "../domain/registration";

/** What the auth feature needs from the outside; tests give a double. */
export interface AuthGateway {
  /** The signed-in user, or null when there is no session (a 401). */
  getMe(): Promise<Me | null>;
  /** Sends the code to the email; the answer is the same whether the account existed or not. */
  register(registration: Registration): Promise<void>;
  /** Sends a login code; an unknown email gets the same answer and no code. */
  login(email: string): Promise<void>;
  /** Opens the session: the back sets the cookie. */
  verifyCode(entry: CodeEntry): Promise<void>;
  /** A new code replaces the previous one. */
  resendCode(email: string): Promise<void>;
  /** Closes the session on the server; it answers the same if there was none. */
  logout(): Promise<void>;
}
