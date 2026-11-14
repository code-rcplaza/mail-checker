import type { FormEventHandler } from "react";
import type { UseFormRegisterReturn } from "react-hook-form";
import { AppLink, type Navigate } from "../../../../components/ui/atoms/AppLink";
import { Button } from "../../../../components/ui/atoms/Button";
import { Input } from "../../../../components/ui/atoms/Input";
import { ForwardIcon } from "../../../../components/ui/atoms/icons";
import { FormField } from "../../../../components/ui/molecules/FormField";
import { t } from "../../../../lib/i18n/translate";

export const registerPath = "/registro";

export interface LoginFormProps {
  readonly fields: { readonly email: UseFormRegisterReturn };
  /** Already translated. */
  readonly errors: Readonly<{ email?: string }>;
  /** What the server said about the whole form (a 429, the network). */
  readonly formError: string | undefined;
  readonly isSubmitting: boolean;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onNavigate: Navigate;
}

export function LoginForm({ fields, errors, formError, isSubmitting, onSubmit, onNavigate }: LoginFormProps) {
  return (
    <div className="flex flex-1 flex-col px-6 pt-8 md:px-0 md:pt-0">
      <div className="flex flex-col gap-2">
        <h1 className="text-[32px] leading-tight font-semibold tracking-[-0.02em] text-ink">
          {t("ui.login.title")}
        </h1>
        <p className="text-[15px] text-ink-secondary">{t("ui.login.subtitle")}</p>
      </div>
      <form noValidate onSubmit={onSubmit} className="flex flex-col gap-6 pt-8">
        <FormField label={t("ui.register.email")} error={errors.email}>
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
        {formError !== undefined && (
          <p role="alert" className="rounded-md bg-error-surface p-4 text-sm text-danger-ink">
            {formError}
          </p>
        )}
        <Button type="submit" loading={isSubmitting}>
          {t("ui.login.submit")}
          <ForwardIcon />
        </Button>
      </form>
      <div className="mt-auto flex flex-col items-center gap-2 pt-6">
        <AppLink
          to={registerPath}
          onNavigate={onNavigate}
          className="inline-flex min-h-11 items-center text-sm font-semibold text-ink-secondary underline"
        >
          {t("ui.login.to_register")}
        </AppLink>
        <p className="text-center text-xs text-ink-muted">{t("ui.login.terms")}</p>
      </div>
    </div>
  );
}
