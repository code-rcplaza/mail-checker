/** Where the person came from to Código: Registro means a first access. */
export const accessFlows = ["register", "login"] as const;

export type AccessFlow = (typeof accessFlows)[number];

export interface CodeEntry {
  readonly email: string;
  readonly code: string;
}

/** With these the code is no longer usable: the screen offers a new one. */
const keysAskingNewCode: readonly string[] = ["auth.code_expired", "auth.code_attempts_exceeded"];

export function asksForNewCode(key: string): boolean {
  return keysAskingNewCode.includes(key);
}
