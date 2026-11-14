import { createFileRoute } from "@tanstack/react-router";
import { RegisterScreen } from "../../features/auth/presentation/screens/RegisterScreen";

export const Route = createFileRoute("/_guest/registro")({ component: RegisterScreen });
