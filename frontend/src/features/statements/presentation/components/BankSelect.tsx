import { FormField } from "../../../../components/ui/molecules/FormField";
import { isMessageKey, t } from "../../../../lib/i18n/translate";
import type { Bank } from "../../domain/upload";

/** A bank the catalog does not name yet still shows, by its id. */
export function bankNameOf(id: string): string {
  const key = `ui.bank.${id}`;
  return isMessageKey(key) ? t(key) : id.toLocaleUpperCase("es-CL");
}

export interface BankSelectProps {
  readonly banks: readonly Bank[];
  /** "" is automatic detection: the upload then sends no `bank`. */
  readonly value: string;
  readonly onChange: (bank: string) => void;
}

export function BankSelect({ banks, value, onChange }: BankSelectProps) {
  return (
    <FormField label={t("ui.upload.bank")}>
      {(control) => (
        <select
          {...control}
          value={value}
          onChange={(event) => onChange(event.target.value)}
          className="min-h-[52px] w-full rounded-md border-[1.5px] border-border-strong bg-field px-4 text-base text-ink focus-visible:border-primary"
        >
          <option value="">{t("ui.upload.bank_auto")}</option>
          {banks.map((bank) => (
            <option key={bank.id} value={bank.id}>
              {bankNameOf(bank.id)}
            </option>
          ))}
        </select>
      )}
    </FormField>
  );
}
