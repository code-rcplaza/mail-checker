import { createFileRoute } from "@tanstack/react-router";
import { RegisterScreen } from "../features/auth/presentation/screens/RegisterScreen";

export const Route = createFileRoute("/registro")({ component: RegisterScreen });
