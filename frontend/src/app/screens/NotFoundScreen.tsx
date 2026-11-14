import { Link } from "@tanstack/react-router";
import { t } from "../../lib/i18n/translate";

export function NotFoundScreen() {
  return (
    <main className="mx-auto flex max-w-md flex-col items-center gap-4 px-6 py-16 text-center">
      <h1 className="text-xl font-semibold text-ink">{t("ui.not_found.title")}</h1>
      <p className="text-ink-secondary">{t("common.route_not_found")}</p>
      <Link to="/" className="inline-flex min-h-11 items-center font-semibold text-primary">
        {t("ui.not_found.back")}
      </Link>
    </main>
  );
}
