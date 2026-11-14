import { describe, expect, it } from "vitest";
import { loginSchema } from "./login";

describe("login schema", () => {
  it.each([
    ["an empty email", { email: "" }, ["auth.email_invalid"]],
    ["an email without format", { email: "ana@" }, ["auth.email_invalid"]],
    ["a valid email", { email: "ana@example.com" }, []],
  ])("reports %s", (_case, input, expected) => {
    const result = loginSchema.safeParse(input);

    expect(result.success ? [] : result.error.issues.map((issue) => issue.message)).toEqual(expected);
  });

  it("trims the email it sends", () => {
    expect(loginSchema.parse({ email: "  ana@example.com " })).toEqual({ email: "ana@example.com" });
  });
});
