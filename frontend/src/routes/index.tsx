import { createFileRoute } from "@tanstack/react-router";

/** Placeholder until the first screens (Bienvenida, Resumen) arrive. */
function HomeScreen() {
  return (
    <main className="mx-auto max-w-md p-6">
      <h1 className="font-wordmark text-2xl font-semibold text-ink">Fokko Subs</h1>
    </main>
  );
}

export const Route = createFileRoute("/")({ component: HomeScreen });
