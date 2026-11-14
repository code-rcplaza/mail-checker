import { useNavigate } from "@tanstack/react-router";
import { AppShell } from "../../../../components/layout/AppShell";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { useNavigateTo } from "../../../../lib/router/useNavigateTo";
import { useThemeToggle } from "../../../../lib/theme/useThemeToggle";
import { useRegister } from "../../application/useRegister";
import { AccessRing } from "../components/AccessRing";
import { RegisterForm } from "../components/RegisterForm";

/** Registro: the container joins the hook, the theme and the router with the form. */
export function RegisterScreen() {
  const navigateTo = useNavigateTo();
  const navigate = useNavigate();
  const theme = useThemeToggle();
  const registration = useRegister({
    onRegistered: (email) => navigate({ to: "/codigo", search: { email, flow: "register" } }),
  });

  return (
    <AppShell header={<BrandHeader theme={theme} />} zone={<AccessRing />}>
      <RegisterForm
        fields={registration.fields}
        errors={registration.errors}
        formError={registration.formError}
        isSubmitting={registration.isSubmitting}
        onSubmit={registration.submit}
        onNavigate={navigateTo}
      />
    </AppShell>
  );
}
