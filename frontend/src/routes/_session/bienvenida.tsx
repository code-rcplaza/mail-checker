import { createFileRoute } from "@tanstack/react-router";
import { WelcomeScreen } from "../../features/auth/presentation/screens/WelcomeScreen";

export const Route = createFileRoute("/_session/bienvenida")({ component: WelcomeScreen });
