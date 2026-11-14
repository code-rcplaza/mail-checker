import { createFileRoute, redirect } from "@tanstack/react-router";
import { z } from "zod";
import { accessFlows } from "../features/auth/domain/code";
import { CodeScreen } from "../features/auth/presentation/screens/CodeScreen";

/** The email the code went to and where the person came from; a broken flow is ignored. */
const codeSearchSchema = z.object({
  email: z.string().optional().catch(undefined),
  flow: z.enum(accessFlows).optional().catch(undefined),
});

function CodeRoute() {
  const { email, flow } = Route.useSearch();
  return <CodeScreen email={email ?? ""} flow={flow} />;
}

export const Route = createFileRoute("/codigo")({
  validateSearch: (search) => codeSearchSchema.parse(search),
  // Opened directly there is no email to check the code against: back to Login.
  beforeLoad: ({ search }) => {
    if (search.email === undefined || search.email === "") {
      throw redirect({ to: "/entrar" });
    }
  },
  component: CodeRoute,
});
