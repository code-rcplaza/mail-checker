import { z } from "zod";

/** The same limit as the back (`statements/adapters/http.DefaultUploadLimits`). */
export const maxUploadBytes = 10 * 1024 * 1024;

/** What the screen knows of a file before sending it. */
export interface FileFacts {
  readonly name: string;
  readonly type: string;
  readonly size: number;
}

function looksLikePDF(file: FileFacts): boolean {
  return file.type === "application/pdf" || (file.type === "" && /\.pdf$/i.test(file.name));
}

/** The i18n key of why the file cannot go, or undefined when it can. Checked again by the back. */
export function fileProblemOf(file: FileFacts): string | undefined {
  if (!looksLikePDF(file)) {
    return "statements.file_kind_unknown";
  }
  if (file.size === 0) {
    return "statements.file_empty";
  }
  if (file.size > maxUploadBytes) {
    return "statements.file_too_large";
  }
  return undefined;
}

const dateInputPattern = /^(\d{4})-(\d{2})-(\d{2})$/;

/**
 * The day of an `<input type="date">` as unix seconds. Noon UTC keeps the same
 * day in Chile and anywhere else, so reading it in the user's calendar never moves it.
 */
export function billingDateToUnix(text: string): number | undefined {
  const [, year, month, day] = dateInputPattern.exec(text) ?? [];
  if (year === undefined || month === undefined || day === undefined) {
    return undefined;
  }
  const noon = new Date(Date.UTC(Number(year), Number(month) - 1, Number(day), 12));
  if (noon.getUTCMonth() !== Number(month) - 1 || noon.getUTCDate() !== Number(day)) {
    return undefined;
  }
  return noon.getTime() / 1000;
}

export interface UploadInput {
  readonly file: File;
  /** Absent: the back detects the bank. */
  readonly bank?: string;
  /** Unix seconds; only when the statement did not bring it. */
  readonly billingDate?: number;
}

/** The multipart of `POST /api/statements`: `file` and only the fields present. */
export function uploadFormOf(input: UploadInput): FormData {
  const form = new FormData();
  form.append("file", input.file);
  if (input.bank !== undefined) {
    form.append("bank", input.bank);
  }
  if (input.billingDate !== undefined) {
    form.append("billingDate", String(input.billingDate));
  }
  return form;
}

/** What `POST /api/statements` answers with a 201. */
export const uploadResultSchema = z.object({
  id: z.string(),
  filename: z.string(),
  months: z.array(z.string()),
  charges: z.number(),
  badLines: z.number(),
});

export type UploadResult = Readonly<z.infer<typeof uploadResultSchema>>;

export const banksSchema = z.array(z.object({ id: z.string() }));

export type Bank = Readonly<z.infer<typeof banksSchema>[number]>;
