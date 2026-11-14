import { useRouter } from "@tanstack/react-router";
import { useCallback } from "react";
import type { Navigate } from "../../components/ui/atoms/AppLink";

/** What screens hand to `AppLink`: a plain path, pushed to the router's history. */
export function useNavigateTo(): Navigate {
  const router = useRouter();
  return useCallback((to: string) => router.history.push(to), [router]);
}
