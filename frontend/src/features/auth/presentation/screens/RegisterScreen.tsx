import { useNavigate, useRouter } from "@tanstack/react-router";
import { AppShell } from "../../../../components/layout/AppShell";
import { BrandHeader } from "../../../../components/layout/BrandHeader";
import { useThemeToggle } from "../../../../lib/theme/useThemeToggle";
import { useRegister } from "../../application/useRegister";
import { AccessRing } from "../components/AccessRing";
import { RegisterForm } from "../components/RegisterForm";

/** Registro: the container joins the hook, the theme and the router with the form. */
export function RegisterScreen() {
  const router = useRouter();
  const navigate = useNavigate();
  const theme = useThemeToggle();
  const registration = useRegister({
    onRegistered: (email) => navigate({ to: "/codigo", search: { email } }),
  });

  return (
    <AppShell header={<BrandHeader theme={theme} />} zone={<AccessRing />}>
      <RegisterForm
        fields={registration.fields}
        errors={registration.errors}
        formError={registration.formError}
        isSubmitting={registration.isSubmitting}
        onSubmit={registration.submit}
        onNavigate={(to) => router.history.push(to)}
      />
    </AppShell>
  );
}
