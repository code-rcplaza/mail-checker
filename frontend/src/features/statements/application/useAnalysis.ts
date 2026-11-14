import { skipToken, useQuery } from "@tanstack/react-query";
import { serverMessageOf } from "../../../lib/api/serverErrors";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import type { UploadResult } from "../domain/upload";
import { foundQueryKey, uploadedStatementKey } from "./useUploadStatement";

/**
 * Análisis: what the upload answered (only in the cache: after a reload it is
 * gone, and the screen goes on without it) and the counts of the summary.
 */
export function useAnalysis(statementId: string) {
  const { statements } = useGateways();
  const uploaded = useQuery<UploadResult>({
    queryKey: uploadedStatementKey(statementId),
    queryFn: skipToken,
  });
  const found = useQuery({ queryKey: foundQueryKey, queryFn: () => statements.readFound() });

  return {
    uploaded: uploaded.data,
    found: found.data,
    isLoading: found.isPending,
    error: found.isError ? serverMessageOf(found.error) : undefined,
  };
}
