import { describe, expect, it } from "vitest";
import { nameMaxLength, registrationSchema } from "./registration";

const valid = { name: "Ana Prueba", email: "ana@example.com", acceptedTerms: true };

/** The i18n key zod reports for each field. */
function errorsOf(input: unknown): Record<string, string> {
  const result = registrationSchema.safeParse(input);
  if (result.success) {
    return {};
  }
  return Object.fromEntries(result.error.issues.map((issue) => [issue.path.join("."), issue.message]));
}

describe("registration schema", () => {
  it.each([
    ["an empty name", { ...valid, name: "" }, { name: "auth.name_required" }],
    ["a name of spaces", { ...valid, name: "   " }, { name: "auth.name_required" }],
    ["a name too long", { ...valid, name: "a".repeat(nameMaxLength + 1) }, { name: "auth.name_too_long" }],
    ["an email without format", { ...valid, email: "ana.example.com" }, { email: "auth.email_invalid" }],
    ["an empty email", { ...valid, email: "" }, { email: "auth.email_invalid" }],
    ["the terms not accepted", { ...valid, acceptedTerms: false }, { acceptedTerms: "auth.terms_required" }],
    ["valid data", valid, {}],
  ])("reports %s", (_case, input, expected) => {
    expect(errorsOf(input)).toEqual(expected);
  });

  it("trims the name and the email it sends", () => {
    expect(registrationSchema.parse({ ...valid, name: "  Ana  ", email: " ana@example.com " })).toEqual({
      ...valid,
      name: "Ana",
      email: "ana@example.com",
    });
  });
});
