import { useNavigate } from "@tanstack/react-router";
import { AppShell } from "../../../../components/layout/AppShell";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { useNavigateTo } from "../../../../lib/router/useNavigateTo";
import { useThemeToggle } from "../../../../lib/theme/useThemeToggle";
import { useLogin } from "../../application/useLogin";
import { AccessRing } from "../components/AccessRing";
import { LoginForm } from "../components/LoginForm";

/** Login: the container joins the hook, the theme and the router with the form. */
export function LoginScreen() {
  const navigate = useNavigate();
  const navigateTo = useNavigateTo();
  const theme = useThemeToggle();
  const login = useLogin({ onCodeSent: (email) => navigate({ to: "/codigo", search: { email } }) });

  return (
    <AppShell header={<BrandHeader theme={theme} />} zone={<AccessRing />}>
      <LoginForm
        fields={login.fields}
        errors={login.errors}
        formError={login.formError}
        isSubmitting={login.isSubmitting}
        onSubmit={login.submit}
        onNavigate={navigateTo}
      />
    </AppShell>
  );
}
