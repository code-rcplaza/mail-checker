import { createFileRoute } from "@tanstack/react-router";
import { LoginScreen } from "../features/auth/presentation/screens/LoginScreen";

export const Route = createFileRoute("/entrar")({ component: LoginScreen });
