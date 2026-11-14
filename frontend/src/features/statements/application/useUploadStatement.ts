import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { ApiError } from "../../../lib/api/client";
import { serverMessageOf } from "../../../lib/api/serverErrors";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import { translateServerKey } from "../../../lib/i18n/translate";
import { billingDateToUnix, fileProblemOf, type UploadInput, type UploadResult } from "../domain/upload";

/** Análisis reads what the upload answered from here. */
export function uploadedStatementKey(id: string) {
  return ["statements", "uploaded", id] as const;
}

const billingDateRequired = "statements.billing_date_required";

export interface UseUploadStatementOptions {
  readonly onUploaded: (result: UploadResult) => void;
}

/** The upload of what was chosen; bank and date go only when present. */
function inputOf(
  file: File,
  choices: { readonly bank: string; readonly billingDate: number | undefined },
): UploadInput {
  return {
    file,
    ...(choices.bank === "" ? {} : { bank: choices.bank }),
    ...(choices.billingDate === undefined ? {} : { billingDate: choices.billingDate }),
  };
}

/** Elegir fuente: a checked PDF, an optional bank and, if the statement lacks it, the billing date. */
export function useUploadStatement({ onUploaded }: UseUploadStatementOptions) {
  const { statements } = useGateways();
  const queryClient = useQueryClient();
  const [file, setFile] = useState<File | undefined>(undefined);
  const [fileError, setFileError] = useState<string | undefined>(undefined);
  const [bank, setBank] = useState("");
  const [billingDateText, setBillingDateText] = useState("");
  const [needsBillingDate, setNeedsBillingDate] = useState(false);
  const [error, setError] = useState<string | undefined>(undefined);

  const uploading = useMutation({
    mutationFn: (input: UploadInput) => statements.upload(input),
    onSuccess: (result) => {
      queryClient.setQueryData(uploadedStatementKey(result.id), result);
      onUploaded(result);
    },
    onError: (caught) => {
      setNeedsBillingDate(
        (already) => already || (caught instanceof ApiError && caught.key === billingDateRequired),
      );
      setError(serverMessageOf(caught));
    },
  });

  function pick(picked: File) {
    const problem = fileProblemOf(picked);
    setFileError(problem === undefined ? undefined : translateServerKey(problem));
    setFile(problem === undefined ? picked : undefined);
    setError(undefined);
  }

  const billingDate = billingDateToUnix(billingDateText);
  const canUpload =
    file !== undefined && !uploading.isPending && (!needsBillingDate || billingDate !== undefined);

  function upload() {
    if (!canUpload || file === undefined) {
      return;
    }
    setError(undefined);
    uploading.mutate(inputOf(file, { bank, billingDate: needsBillingDate ? billingDate : undefined }));
  }

  return {
    file,
    fileError,
    pick,
    bank,
    chooseBank: setBank,
    needsBillingDate,
    billingDateText,
    writeBillingDate: setBillingDateText,
    error,
    isUploading: uploading.isPending,
    canUpload,
    upload,
  };
}
