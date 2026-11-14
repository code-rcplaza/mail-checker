import { createFileRoute } from "@tanstack/react-router";
import { UploadScreen } from "../../features/statements/presentation/screens/UploadScreen";

function SourceRoute() {
  const { me } = Route.useRouteContext();
  return <UploadScreen account={{ name: me.name, email: me.email }} />;
}

export const Route = createFileRoute("/_session/fuente")({ component: SourceRoute });
