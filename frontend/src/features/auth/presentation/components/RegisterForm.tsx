import type { FormEventHandler } from "react";
import type { UseFormRegisterReturn } from "react-hook-form";
import { AppLink, type Navigate } from "../../../../components/ui/atoms/AppLink";
import { Button } from "../../../../components/ui/atoms/Button";
import { Checkbox } from "../../../../components/ui/atoms/Checkbox";
import { Input } from "../../../../components/ui/atoms/Input";
import { ForwardIcon, ShieldIcon } from "../../../../components/ui/atoms/icons";
import { FormField } from "../../../../components/ui/molecules/FormField";
import { t } from "../../../../lib/i18n/translate";

/** Where Login lives; the terms link waits for the legal pages (FOKKO-0113A-SUBS). */
export const loginPath = "/entrar";

export interface RegisterFormProps {
  readonly fields: {
    readonly name: UseFormRegisterReturn;
    readonly email: UseFormRegisterReturn;
    readonly acceptedTerms: UseFormRegisterReturn;
  };
  /** Already translated, one per field. */
  readonly errors: Readonly<Partial<Record<"name" | "email" | "acceptedTerms", string>>>;
  /** What the server said about the whole form (a 429, the network). */
  readonly formError: string | undefined;
  readonly isSubmitting: boolean;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onNavigate: Navigate;
}

export function RegisterForm({
  fields,
  errors,
  formError,
  isSubmitting,
  onSubmit,
  onNavigate,
}: RegisterFormProps) {
  return (
    <div className="flex flex-1 flex-col px-6 pt-8 md:px-0 md:pt-0">
      <div className="flex flex-col gap-2">
        <h1 className="text-[32px] leading-tight font-semibold tracking-[-0.02em] text-ink">
          {t("ui.register.title")}
        </h1>
        <p className="text-[15px] text-ink-secondary">{t("ui.register.subtitle")}</p>
      </div>
      <form noValidate onSubmit={onSubmit} className="flex flex-col gap-6 pt-8">
        <FormField label={t("ui.register.name")} error={errors.name}>
          {(control) => (
            <Input
              {...control}
              {...fields.name}
              type="text"
              autoComplete="name"
              placeholder={t("ui.register.name_placeholder")}
            />
          )}
        </FormField>
        <FormField label={t("ui.register.email")} help={t("ui.register.email_help")} error={errors.email}>
          {(control) => (
            <Input
              {...control}
              {...fields.email}
              type="email"
              inputMode="email"
              autoComplete="email"
              placeholder={t("ui.register.email_placeholder")}
            />
          )}
        </FormField>
        <p className="flex items-start gap-4 text-sm leading-normal text-ink-secondary">
          <span className="mt-1 flex flex-none text-primary">
            <ShieldIcon />
          </span>
          {t("ui.register.privacy")}
        </p>
        <Checkbox {...fields.acceptedTerms} label={t("ui.register.terms")} error={errors.acceptedTerms} />
        {formError !== undefined && (
          <p role="alert" className="rounded-md bg-error-surface p-4 text-sm text-danger-ink">
            {formError}
          </p>
        )}
        <Button type="submit" loading={isSubmitting}>
          {t("ui.register.submit")}
          <ForwardIcon />
        </Button>
      </form>
      <div className="mt-auto flex justify-center pt-6">
        <AppLink
          to={loginPath}
          onNavigate={onNavigate}
          className="inline-flex min-h-11 items-center text-sm font-semibold text-ink-secondary underline"
        >
          {t("ui.register.to_login")}
        </AppLink>
      </div>
    </div>
  );
}
