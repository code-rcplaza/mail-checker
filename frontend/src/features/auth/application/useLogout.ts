import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useGateways } from "../../../lib/gateways/GatewaysProvider";

export interface UseLogoutOptions {
  /** `reachedServer` is false when the server could not confirm the close (offline). */
  readonly onLoggedOut: (reachedServer: boolean) => void;
}

/** "Cerrar sesión": whatever the server says, nothing of the account stays in the cache. */
export function useLogout({ onLoggedOut }: UseLogoutOptions) {
  const { auth } = useGateways();
  const queryClient = useQueryClient();
  const leaving = useMutation({
    mutationFn: () => auth.logout(),
    onSettled: (_answer, error) => {
      queryClient.removeQueries();
      onLoggedOut(error === null);
    },
  });

  return { logout: () => leaving.mutate(), isLeaving: leaving.isPending };
}
