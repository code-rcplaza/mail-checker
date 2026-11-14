import type { Bank, UploadInput, UploadResult } from "../domain/upload";

/** What the statements feature needs from the outside; tests give a double. */
export interface StatementsGateway {
  listBanks(): Promise<readonly Bank[]>;
  upload(input: UploadInput): Promise<UploadResult>;
}
