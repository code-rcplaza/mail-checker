import { createFileRoute } from "@tanstack/react-router";
import { AnalysisScreen } from "../../../features/statements/presentation/screens/AnalysisScreen";

function AnalysisRoute() {
  const { id } = Route.useParams();
  const { me } = Route.useRouteContext();
  return <AnalysisScreen statementId={id} account={{ name: me.name, email: me.email }} />;
}

export const Route = createFileRoute("/_session/analisis/$id")({ component: AnalysisRoute });
