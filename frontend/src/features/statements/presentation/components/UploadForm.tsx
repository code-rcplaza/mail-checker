import { Button } from "../../../../components/ui/atoms/Button";
import { Input } from "../../../../components/ui/atoms/Input";
import { ShieldIcon } from "../../../../components/ui/atoms/icons";
import { FormField } from "../../../../components/ui/molecules/FormField";
import { t } from "../../../../lib/i18n/translate";
import type { Bank } from "../../domain/upload";
import { BankSelect } from "./BankSelect";
import { UploadDropzone } from "./UploadDropzone";

export interface UploadFormProps {
  readonly file: File | undefined;
  readonly fileError: string | undefined;
  readonly onPick: (file: File) => void;
  readonly banks: readonly Bank[];
  readonly bank: string;
  readonly onChooseBank: (bank: string) => void;
  /** Only after the back said the statement brings no billing date. */
  readonly needsBillingDate: boolean;
  readonly billingDate: string;
  readonly onWriteBillingDate: (date: string) => void;
  /** Already translated: why the upload failed. */
  readonly error: string | undefined;
  readonly canUpload: boolean;
  readonly isUploading: boolean;
  readonly onUpload: () => void;
}

export function UploadForm(props: UploadFormProps) {
  return (
    <div className="flex flex-col gap-8 px-6 pt-8 md:px-0 md:pt-0">
      <section className="flex flex-col gap-4">
        <h2 className="text-[17px] font-semibold tracking-[-0.01em] text-ink">{t("ui.upload.section")}</h2>
        <p className="text-sm text-ink-secondary">{t("ui.upload.section_text")}</p>
        <UploadDropzone file={props.file} error={props.fileError} onPick={props.onPick} />
        <BankSelect banks={props.banks} value={props.bank} onChange={props.onChooseBank} />
        {props.needsBillingDate && (
          <FormField label={t("ui.upload.billing_date")}>
            {(control) => (
              <Input
                {...control}
                type="date"
                value={props.billingDate}
                onChange={(event) => props.onWriteBillingDate(event.target.value)}
              />
            )}
          </FormField>
        )}
        {props.error !== undefined && (
          <p role="alert" className="rounded-md bg-error-surface p-4 text-sm text-danger-ink">
            {props.error}
          </p>
        )}
        <Button disabled={!props.canUpload} loading={props.isUploading} onClick={props.onUpload}>
          {t("ui.upload.submit")}
        </Button>
      </section>
      <section className="flex items-start gap-4 border-t border-border pt-8">
        <span className="mt-1 flex flex-none text-primary">
          <ShieldIcon />
        </span>
        <div className="flex flex-col gap-1">
          <h2 className="text-[15px] font-semibold text-ink">{t("ui.upload.kept_title")}</h2>
          <p className="text-[13px] text-ink-secondary">{t("ui.upload.kept_text")}</p>
        </div>
      </section>
    </div>
  );
}
