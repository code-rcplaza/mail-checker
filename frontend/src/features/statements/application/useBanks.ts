import { useQuery } from "@tanstack/react-query";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";
import type { Bank } from "../domain/upload";

const noBanks: readonly Bank[] = [];

/** The bank list is a help: if it fails, the select stays on automatic and the upload still goes. */
export function useBanks(): readonly Bank[] {
  const { statements } = useGateways();
  const reading = useQuery({
    queryKey: ["banks"],
    queryFn: () => statements.listBanks(),
    staleTime: Number.POSITIVE_INFINITY,
  });
  return reading.data ?? noBanks;
}
