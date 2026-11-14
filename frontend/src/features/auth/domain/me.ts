import { z } from "zod";

/** The signed-in user, as `GET /api/me` sends it. */
export const meSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string(),
  initials: z.string(),
});

export type Me = Readonly<z.infer<typeof meSchema>>;

/** The greeting uses only the first word of the name. */
export function firstNameOf(name: string): string {
  return name.trim().split(/\s+/)[0] ?? "";
}
