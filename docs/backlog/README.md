# Backlog

Loop: `/idea` → `/stories NNN` → `/story NNN-MM` (TDD). Reglas en `CLAUDE.md`.

- Cada idea vive en `NNN-slug/` con su `idea.md` y una historia por archivo `MM-slug.md`.
- El estado de las ideas vive **solo** en esta tabla; el de las historias, **solo** en
  el checklist de su `idea.md`.
- Las ideas terminadas no se mueven ni se borran.
- Este backlog es historial; el diseño vigente está en `docs/ARCHITECTURE.md`.

| # | Idea | Estado |
|---|---|---|
| FOKKO-0001-SUBS | [Armazón del backend (fase 1)](FOKKO-0001-SUBS-foundation/idea.md) | done |
| FOKKO-0002-SUBS | [Núcleo de suscripciones (fase 2, parte 1)](FOKKO-0002-SUBS-subscriptions-core/idea.md) | done |
| FOKKO-0003-SUBS | [Núcleo de sincronización de correos (fase 2, parte 2)](FOKKO-0003-SUBS-emailsync-core/idea.md) | done |
| FOKKO-0004-SUBS | [Núcleo de cartolas (fase 2, parte 3)](FOKKO-0004-SUBS-statements-core/idea.md) | done |
| FOKKO-0005-SUBS | [Núcleo de cuentas (fase 2, parte 4)](FOKKO-0005-SUBS-accounts-core/idea.md) | done |
| FOKKO-0006-SUBS | [Catálogo de servicios (fase 3, parte 1)](FOKKO-0006-SUBS-service-catalog/idea.md) | done |
| FOKKO-0007-SUBS | [Líneas de suscripción (corrección de la detección)](FOKKO-0007-SUBS-subscription-lines/idea.md) | done |

## Plan del MVP (ideas por hacer)

Cada idea es unitaria: un resultado verificable, una capa, una feature, ~300 líneas. Las filas
sin enlace aún no tienen `idea.md`; se crea con `/idea` cuando les toca. El orden numérico
respeta las dependencias. ⛔ marca una decisión abierta que la bloquea.

### A. Catálogo y cartolas

| # | Idea | Estado |
|---|---|---|
| FOKKO-0008-SUBS | [Enlace de cancelación en el catálogo](FOKKO-0008-SUBS-cancel-link/idea.md) | doing |
| FOKKO-0009-SUBS | [`extract`: texto de un PDF](FOKKO-0009-SUBS-pdf-extract/idea.md) | done |
| FOKKO-0010-SUBS | [Meses leídos de una cartola](FOKKO-0010-SUBS-months-read/idea.md) | done |
| FOKKO-0011A-SUBS | [Parser BCI: reconocer la cartola](FOKKO-0011A-SUBS-bci-detect/idea.md) | done |
| FOKKO-0011B-SUBS | [Parser BCI: leer los cobros del estado nacional](FOKKO-0011B-SUBS-bci-national-lines/idea.md) | done |
| FOKKO-0012-SUBS | [Decidir si un cobro internacional es en pesos](FOKKO-0012-SUBS-pesos-or-usd/idea.md) | done |
| FOKKO-0013-SUBS | [Parser BCI: cobros del estado internacional (origen y US$)](FOKKO-0013-SUBS-bci-international/idea.md) | done |
| FOKKO-0014-SUBS | [Parser BCI: par de traspasos](FOKKO-0014-SUBS-bci-transfer-pair/idea.md) | done |
| FOKKO-0015-SUBS | [Derivar el tipo de cambio de los traspasos](FOKKO-0015-SUBS-derive-rate/idea.md) | done |
| FOKKO-0016-SUBS | [Último tipo de cambio o total incompleto](FOKKO-0016-SUBS-latest-rate/idea.md) | done |
| FOKKO-0017-SUBS | [Convertir un cobro en US$ a CLP](FOKKO-0017-SUBS-convert-usd/idea.md) | done |
| FOKKO-0018A-SUBS | [Puerto `BillingDate` en `BankStatementParser`](FOKKO-0018A-SUBS-billing-date-port/idea.md) | done |
| FOKKO-0018B-SUBS | [`ImportStatement` toma la fecha del archivo](FOKKO-0018B-SUBS-import-file-date/idea.md) | done |
| FOKKO-0019-SUBS | [Parser BCI: fecha de facturación](FOKKO-0019-SUBS-bci-billing-date/idea.md) | done |

### B. Detección

| # | Idea | Estado |
|---|---|---|
| FOKKO-0020A-SUBS | [`UserID` en `Charge`](FOKKO-0020A-SUBS-charge-user-id/idea.md) | done |
| FOKKO-0020B-SUBS | [`statements` filtra por usuario](FOKKO-0020B-SUBS-statements-by-user/idea.md) | done |
| FOKKO-0021-SUBS | [`UserID` en los puertos de `subscriptions`](FOKKO-0021-SUBS-user-id-subscriptions/idea.md) | done |
| FOKKO-0022-SUBS | [Deshacer una decisión](FOKKO-0022-SUBS-undo-decision/idea.md) | done |
| FOKKO-0023-SUBS | [Próximo cobro por línea](FOKKO-0023-SUBS-next-charge/idea.md) | done |
| FOKKO-0024-SUBS | [Cobros ya cobrados del mes](FOKKO-0024-SUBS-charged-this-month/idea.md) | done |
| FOKKO-0025-SUBS | [Cobros restantes del mes](FOKKO-0025-SUBS-remaining-this-month/idea.md) | done |
| FOKKO-0026-SUBS | [Detectar un alza](FOKKO-0026-SUBS-detect-rise/idea.md) | done |
| FOKKO-0027-SUBS | [Impacto mensual y anual del alza](FOKKO-0027-SUBS-rise-impact/idea.md) | done |
| FOKKO-0028-SUBS | [Descartar el aviso de alza](FOKKO-0028-SUBS-dismiss-rise-notice/idea.md) | done |
| FOKKO-0029-SUBS | [Total del mes en CLP](FOKKO-0029-SUBS-total-clp/idea.md) | done |
| FOKKO-0030-SUBS | [Monto original junto al equivalente ("≈")](FOKKO-0030-SUBS-original-amount/idea.md) | done |
| FOKKO-0031-SUBS | [Crear suscripción manual](FOKKO-0031-SUBS-create-manual/idea.md) | done |
| FOKKO-0032-SUBS | [Editar o quitar una manual](FOKKO-0032-SUBS-edit-manual/idea.md) | done |
| FOKKO-0033-SUBS | [Override de monto y día](FOKKO-0033-SUBS-line-override/idea.md) | done |
| FOKKO-0034-SUBS | ["Dejar de seguir" una línea](FOKKO-0034-SUBS-stop-following/idea.md) | done |
| FOKKO-0035-SUBS | [Marcar "por cancelar"](FOKKO-0035-SUBS-mark-for-cancel/idea.md) | done |
| FOKKO-0036-SUBS | [`Summary`: reparto por servicio](FOKKO-0036-SUBS-summary-spend/idea.md) | done |
| FOKKO-0037-SUBS | [`Summary`: acciones pendientes](FOKKO-0037-SUBS-pending-actions/idea.md) | done |
| FOKKO-0038-SUBS | [`Summary`: próximos cobros](FOKKO-0038-SUBS-upcoming-charges/idea.md) | done |
| FOKKO-0039-SUBS | [`Summary`: detalle de una línea](FOKKO-0039-SUBS-line-detail/idea.md) | done |

### C. Acceso

| # | Idea | Estado |
|---|---|---|
| FOKKO-0040-SUBS | [Registrar usuario](FOKKO-0040-SUBS-auth-user/idea.md) | done |
| FOKKO-0041-SUBS | [Pedir código al entrar (login)](FOKKO-0041-SUBS-auth-login/idea.md) | done |
| FOKKO-0042-SUBS | [Generar el código](FOKKO-0042-SUBS-auth-issue-code/idea.md) | done |
| FOKKO-0043-SUBS | [Verificar el código (vigencia y 5 intentos)](FOKKO-0043-SUBS-auth-verify-code/idea.md) | done |
| FOKKO-0044-SUBS | [Límite de 3 envíos por hora](FOKKO-0044-SUBS-auth-send-limit/idea.md) | done |
| FOKKO-0045-SUBS | [Puerto de envío con doble de consola](FOKKO-0045-SUBS-auth-code-sender/idea.md) | done |
| FOKKO-0046-SUBS | [Crear sesión de 30 días](FOKKO-0046-SUBS-auth-open-session/idea.md) | done |
| FOKKO-0047A-SUBS | [Validar la sesión](FOKKO-0047A-SUBS-validate-session/idea.md) | done |
| FOKKO-0047B-SUBS | [Cerrar la sesión](FOKKO-0047B-SUBS-close-session/idea.md) | done |
| FOKKO-0048-SUBS | [Envío del código con Resend (`net/http`, timeout, llave por entorno)](FOKKO-0048-SUBS-auth-resend-mailer/idea.md) | done |

### D. Persistencia

| # | Idea | Estado |
|---|---|---|
| FOKKO-0049-SUBS | [`sqlc.yaml` y generación de queries (la migración base no crea tablas)](FOKKO-0049-SUBS-sqlc-setup/idea.md) | done |
| FOKKO-0050-SUBS | [Repo de usuarios](FOKKO-0050-SUBS-repo-users/idea.md) | done |
| FOKKO-0051A-SUBS | [Repo de códigos de acceso](FOKKO-0051A-SUBS-repo-login-codes/idea.md) | done |
| FOKKO-0051B-SUBS | [Repo del registro de envíos de códigos](FOKKO-0051B-SUBS-repo-code-sends/idea.md) | done |
| FOKKO-0052-SUBS | [Repo de sesiones](FOKKO-0052-SUBS-repo-sessions/idea.md) | done |
| FOKKO-0053A-SUBS | [Repo de cartolas: crear con sus cargos](FOKKO-0053A-SUBS-repo-statements-create/idea.md) | done |
| FOKKO-0053B-SUBS | [Repo de cartolas: lectura y borrado](FOKKO-0053B-SUBS-repo-statements-read-delete/idea.md) | done |
| FOKKO-0053C-SUBS | [Lector de cargos de `subscriptions`](FOKKO-0053C-SUBS-charge-reader/idea.md) | done |
| FOKKO-0054-SUBS | [Repo de decisiones](FOKKO-0054-SUBS-repo-decisions/idea.md) | done |
| FOKKO-0055-SUBS | [Repo de alias](FOKKO-0055-SUBS-repo-aliases/idea.md) | done |
| FOKKO-0056-SUBS | [Repo de overrides](FOKKO-0056-SUBS-repo-overrides/idea.md) | done |
| FOKKO-0057-SUBS | [Repo de manuales](FOKKO-0057-SUBS-repo-manuals/idea.md) | done |
| FOKKO-0058A-SUBS | [Repo de líneas ocultas](FOKKO-0058A-SUBS-repo-hidden-lines/idea.md) | done |
| FOKKO-0058B-SUBS | [Repo de cancelaciones](FOKKO-0058B-SUBS-repo-cancellations/idea.md) | done |
| FOKKO-0059-SUBS | [Repo de decisiones de alza](FOKKO-0059-SUBS-repo-price-alerts/idea.md) | done |
| FOKKO-0060-SUBS | [Repo del tipo de cambio](FOKKO-0060-SUBS-repo-exchange-rates/idea.md) | done |
| FOKKO-0061-SUBS | [Limpiar códigos y sesiones vencidos](FOKKO-0061-SUBS-cleanup-expired/idea.md) | done |

### E. HTTP

| # | Idea | Estado |
|---|---|---|
| FOKKO-0062A-SUBS | [Recuperar panics](FOKKO-0062A-SUBS-panic-recovery/idea.md) | done |
| FOKKO-0062B-SUBS | [Log de peticiones sin datos personales](FOKKO-0062B-SUBS-request-logging/idea.md) | done |
| FOKKO-0062C-SUBS | [Timeouts del servidor](FOKKO-0062C-SUBS-server-timeouts/idea.md) | done |
| FOKKO-0063-SUBS | [Errores con llaves i18n](FOKKO-0063-SUBS-http-errors/idea.md) | done |
| FOKKO-0064-SUBS | [Middleware de sesión](FOKKO-0064-SUBS-http-session-middleware/idea.md) | done |
| FOKKO-0065A-SUBS | [Cookie de sesión](FOKKO-0065A-SUBS-session-cookie/idea.md) | done |
| FOKKO-0065B-SUBS | [CSRF: comprobar el origen](FOKKO-0065B-SUBS-csrf-origin/idea.md) | done |
| FOKKO-0065C-SUBS | [Sin cabeceras CORS](FOKKO-0065C-SUBS-no-cors/idea.md) | done |
| FOKKO-0066A-SUBS | [`POST /api/auth/register`](FOKKO-0066A-SUBS-register-endpoint/idea.md) | done |
| FOKKO-0066B-SUBS | [`POST /api/auth/login`](FOKKO-0066B-SUBS-login-endpoint/idea.md) | done |
| FOKKO-0067A-SUBS | [`POST /api/auth/code/verify`](FOKKO-0067A-SUBS-verify-endpoint/idea.md) | done |
| FOKKO-0067B-SUBS | [`POST /api/auth/code/resend`](FOKKO-0067B-SUBS-resend-endpoint/idea.md) | done |
| FOKKO-0068A-SUBS | [`GET /api/me`](FOKKO-0068A-SUBS-me-endpoint/idea.md) | done |
| FOKKO-0068B-SUBS | [`POST /api/auth/logout`](FOKKO-0068B-SUBS-logout-endpoint/idea.md) | done |
| FOKKO-0069-SUBS | [Límites de la subida (tamaño, PDF, timeout)](FOKKO-0069-SUBS-http-upload-limits/idea.md) | done |
| FOKKO-0070-SUBS | [`POST /api/statements`](FOKKO-0070-SUBS-http-statements-upload/idea.md) | done |
| FOKKO-0071-SUBS | [`GET /api/banks`](FOKKO-0071-SUBS-http-banks/idea.md) | done |
| FOKKO-0072A-SUBS | [`GET /api/summary`: parámetro `months`](FOKKO-0072A-SUBS-summary-months-param/idea.md) | done |
| FOKKO-0072B-SUBS | [`GET /api/summary`: suscripciones](FOKKO-0072B-SUBS-summary-subscriptions/idea.md) | done |
| FOKKO-0072C-SUBS | [`GET /api/summary`: posibles con su evidencia](FOKKO-0072C-SUBS-summary-possible/idea.md) | done |
| FOKKO-0072D-SUBS | [`GET /api/summary`: total del mes](FOKKO-0072D-SUBS-summary-total/idea.md) | done |
| FOKKO-0072E-SUBS | [`GET /api/summary`: cobrados y restantes del mes](FOKKO-0072E-SUBS-summary-month-counts/idea.md) | done |
| FOKKO-0072F-SUBS | [`GET /api/summary`: alzas pendientes](FOKKO-0072F-SUBS-summary-price-rises/idea.md) | done |
| FOKKO-0072G-SUBS | [`GET /api/summary`: acciones pendientes](FOKKO-0072G-SUBS-summary-pending-actions/idea.md) | done |
| FOKKO-0072H-SUBS | [`GET /api/summary`: próximos cobros](FOKKO-0072H-SUBS-summary-upcoming/idea.md) | done |
| FOKKO-0073A-SUBS | [`POST /api/subscriptions/decisions`](FOKKO-0073A-SUBS-decision-post/idea.md) | done |
| FOKKO-0073B-SUBS | [`DELETE /api/subscriptions/decisions`](FOKKO-0073B-SUBS-decision-delete/idea.md) | done |
| FOKKO-0074A-SUBS | [`PUT /api/subscriptions/aliases`](FOKKO-0074A-SUBS-alias-put/idea.md) | done |
| FOKKO-0074B-SUBS | [`PUT /api/subscriptions/overrides`](FOKKO-0074B-SUBS-override-put/idea.md) | done |
| FOKKO-0075-SUBS | [Manuales](FOKKO-0075-SUBS-http-manuals/idea.md) | done |
| FOKKO-0076-SUBS | [Descartar aviso de alza](FOKKO-0076-SUBS-http-price-alert/idea.md) | done |
| FOKKO-0077A-SUBS | [Guía de cancelación: caso de uso y `GET`](FOKKO-0077A-SUBS-cancel-guide/idea.md) | done |
| FOKKO-0077B-SUBS | [`POST /api/subscriptions/cancellations`](FOKKO-0077B-SUBS-cancellation-post/idea.md) | done |
| FOKKO-0078-SUBS | [`DELETE /api/subscriptions/line`](FOKKO-0078-SUBS-http-stop-following/idea.md) | done |
| FOKKO-0079-SUBS | [Exportar copia (CSV/JSON)](FOKKO-0079-SUBS-http-export/idea.md) | done |
| FOKKO-0080-SUBS | [`DELETE /api/me`](FOKKO-0080-SUBS-http-delete-me/idea.md) | done |
| FOKKO-0081-SUBS | [`GET /api/health` también revisa la base (para el contenedor)](FOKKO-0081-SUBS-http-health-db/idea.md) | done |
| FOKKO-0082A-SUBS | [Armar las dependencias](FOKKO-0082A-SUBS-wire-dependencies/idea.md) | done |
| FOKKO-0082B-SUBS | [Registrar las rutas](FOKKO-0082B-SUBS-wire-routes/idea.md) | done |
| FOKKO-0082C-SUBS | [Limpieza periódica](FOKKO-0082C-SUBS-wire-cleanup/idea.md) | done |
| FOKKO-0082D-SUBS | [Arranque que falla rápido](FOKKO-0082D-SUBS-fail-fast-startup/idea.md) | done |

### F. Front

| # | Idea | Estado |
|---|---|---|
| FOKKO-0083A-SUBS | [Proyecto Vite con React](FOKKO-0083A-SUBS-vite-react/idea.md) | done |
| FOKKO-0083B-SUBS | [TypeScript estricto con Biome](FOKKO-0083B-SUBS-strict-biome/idea.md) | done |
| FOKKO-0083C-SUBS | [Corredor de pruebas del front](FOKKO-0083C-SUBS-front-tests/idea.md) | done |
| FOKKO-0084A-SUBS | [Tokens de diseño en Tailwind](FOKKO-0084A-SUBS-tailwind-tokens/idea.md) | done |
| FOKKO-0084B-SUBS | [Preferencia de tema](FOKKO-0084B-SUBS-theme-preference/idea.md) | done |
| FOKKO-0085A-SUBS | [Cliente de la API](FOKKO-0085A-SUBS-api-client/idea.md) | done |
| FOKKO-0085B-SUBS | [Cimientos de TanStack (Query, Router)](FOKKO-0085B-SUBS-query-router-base/idea.md) | done |
| FOKKO-0085C-SUBS | [Límites entre capas del front](FOKKO-0085C-SUBS-front-layer-boundaries/idea.md) | done |
| FOKKO-0085D-SUBS | [Puertos del front por Context](FOKKO-0085D-SUBS-front-gateways-context/idea.md) | done |
| FOKKO-0086A-SUBS | [Traducción tipada con catálogo en español](FOKKO-0086A-SUBS-translate-catalog/idea.md) | done |
| FOKKO-0086B-SUBS | [Test de contrato con las llaves del back](FOKKO-0086B-SUBS-i18n-contract/idea.md) | done |
| FOKKO-0087-SUBS | [`Button`, `Input`, `FormField`, `Checkbox`](FOKKO-0087-SUBS-front-form-components/idea.md) | done |
| FOKKO-0088-SUBS | [`Card`, `Badge`, `ListRow`, `ServiceAvatar`](FOKKO-0088-SUBS-front-display-components/idea.md) | done |
| FOKKO-0089A-SUBS | [`SegmentedRing`](FOKKO-0089A-SUBS-segmented-ring/idea.md) | done |
| FOKKO-0089B-SUBS | [`AnimatedNumber`](FOKKO-0089B-SUBS-animated-number/idea.md) | done |
| FOKKO-0090-SUBS | [`Toast`, `Skeleton`, `SegmentedControl`, `EmptyState`](FOKKO-0090-SUBS-front-feedback-components/idea.md) | done |
| FOKKO-0091-SUBS | [`CodeInput`](FOKKO-0091-SUBS-front-code-input/idea.md) | done |
| FOKKO-0092-SUBS | [`AppShell` móvil](FOKKO-0092-SUBS-front-shell-mobile/idea.md) | done |
| FOKKO-0093-SUBS | [`SideNav` de escritorio](FOKKO-0093-SUBS-front-shell-desktop/idea.md) | done |
| FOKKO-0094-SUBS | [Pantalla Registro](FOKKO-0094-SUBS-screen-register/idea.md) | doing |
| FOKKO-0095-SUBS | [Pantalla Login](FOKKO-0095-SUBS-screen-login/idea.md) | doing |
| FOKKO-0096-SUBS | [Pantalla Código](FOKKO-0096-SUBS-screen-code/idea.md) | todo |
| FOKKO-0097-SUBS | [Pantalla Bienvenida](FOKKO-0097-SUBS-screen-welcome/idea.md) | todo |
| FOKKO-0098-SUBS | [Protección de rutas](FOKKO-0098-SUBS-front-session-routes/idea.md) | todo |
| FOKKO-0099-SUBS | [Subir cartola](FOKKO-0099-SUBS-screen-upload/idea.md) | todo |
| FOKKO-0100A-SUBS | [Pantalla Cargando](FOKKO-0100A-SUBS-screen-loading/idea.md) | todo |
| FOKKO-0100B-SUBS | [Pantalla Análisis](FOKKO-0100B-SUBS-screen-analysis-result/idea.md) | todo |
| FOKKO-0101-SUBS | [Error de lectura](FOKKO-0101-SUBS-screen-upload-error/idea.md) | todo |
| FOKKO-0102-SUBS | [Revisión (con selector de cadencia)](FOKKO-0102-SUBS-screen-review/idea.md) | todo |
| FOKKO-0103-SUBS | [Resumen](FOKKO-0103-SUBS-screen-summary/idea.md) | todo |
| FOKKO-0104-SUBS | [`CycleTip` con el texto corregido](FOKKO-0104-SUBS-screen-cycle-tip/idea.md) | todo |
| FOKKO-0105-SUBS | [Subs](FOKKO-0105-SUBS-screen-subs/idea.md) | todo |
| FOKKO-0106A-SUBS | [Detalle de una línea (con o sin alza)](FOKKO-0106A-SUBS-detail-line/idea.md) | todo |
| FOKKO-0106B-SUBS | [Dejar de seguir](FOKKO-0106B-SUBS-stop-following-ui/idea.md) | todo |
| FOKKO-0107A-SUBS | [Agregar una suscripción manual](FOKKO-0107A-SUBS-add-manual/idea.md) | todo |
| FOKKO-0107B-SUBS | [Editar una suscripción](FOKKO-0107B-SUBS-edit-subscription/idea.md) | todo |
| FOKKO-0108-SUBS | [Alerta de alza](FOKKO-0108-SUBS-screen-price-alert/idea.md) | todo |
| FOKKO-0109-SUBS | [Cómo cancelar](FOKKO-0109-SUBS-screen-cancel-guide/idea.md) | todo |
| FOKKO-0110A-SUBS | [Perfil](FOKKO-0110A-SUBS-profile/idea.md) | todo |
| FOKKO-0110B-SUBS | [Apariencia](FOKKO-0110B-SUBS-appearance/idea.md) | todo |
| FOKKO-0111-SUBS | [Eliminar todo](FOKKO-0111-SUBS-screen-delete-all/idea.md) | todo |

### G. Cierre

| # | Idea | Estado |
|---|---|---|
| FOKKO-0112A-SUBS | [Smoke test: acceso](FOKKO-0112A-SUBS-smoke-auth/idea.md) | done |
| FOKKO-0112B-SUBS | [Smoke test: de la cartola al resumen](FOKKO-0112B-SUBS-smoke-import-summary/idea.md) | todo |
| FOKKO-0112C-SUBS | [Smoke test: aislamiento entre usuarios](FOKKO-0112C-SUBS-smoke-isolation/idea.md) | done |
| FOKKO-0112D-SUBS | [Smoke test: eliminar todo](FOKKO-0112D-SUBS-smoke-delete-all/idea.md) | done |
| FOKKO-0113A-SUBS | [Textos legales ⛔ decisión 12](FOKKO-0113A-SUBS-legal-pages/idea.md) | todo |
| FOKKO-0113B-SUBS | [Pruebas de privacidad](FOKKO-0113B-SUBS-privacy-tests/idea.md) | todo |
| FOKKO-0114A-SUBS | [Embeber el front en el binario](FOKKO-0114A-SUBS-embed-front/idea.md) | todo |
| FOKKO-0114B-SUBS | [Respaldo de rutas del front](FOKKO-0114B-SUBS-spa-fallback/idea.md) | todo |
| FOKKO-0114C-SUBS | [Cabeceras de caché](FOKKO-0114C-SUBS-static-cache/idea.md) | todo |
| FOKKO-0114D-SUBS | [Comando de build único](FOKKO-0114D-SUBS-build-command/idea.md) | todo |
| FOKKO-0115-SUBS | [Dockerfile multi-stage](FOKKO-0115-SUBS-dockerfile/idea.md) | todo |
| FOKKO-0116-SUBS | [`docker-compose` local](FOKKO-0116-SUBS-docker-compose-local/idea.md) | todo |
| FOKKO-0117-SUBS | [README](FOKKO-0117-SUBS-readme/idea.md) | todo |

### H. Agregadas por la auditoría del código

| # | Idea | Estado |
|---|---|---|
| FOKKO-0118-SUBS | [Guardar el tipo de cambio derivado al importar una cartola](FOKKO-0118-SUBS-record-exchange-rate/idea.md) | done |

### Post-MVP (sin número aún)

Lectura de correo por IMAP (cuentas, sincronización, avisos BCI, conciliación), `later` de la
alerta de alza, recordatorio "2 días antes", CSV, OAuth de Gmail, login con Google y Apple,
calendario, notificaciones, recuperar acceso y uso sin conexión.
