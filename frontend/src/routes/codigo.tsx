import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";

/** The email the code went to; Registro and Login send it here. */
const codeSearchSchema = z.object({ email: z.string().optional() });

/** Placeholder until the Código screen (FOKKO-0096-SUBS). */
function CodeScreen() {
  return <main />;
}

export const Route = createFileRoute("/codigo")({
  validateSearch: (search) => codeSearchSchema.parse(search),
  component: CodeScreen,
});
