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
| FOKKO-0006-SUBS | [Catálogo de servicios (fase 3, parte 1)](FOKKO-0006-SUBS-service-catalog/idea.md) | doing |
| FOKKO-0007-SUBS | [Líneas de suscripción (corrección de la detección)](FOKKO-0007-SUBS-subscription-lines/idea.md) | done |

## Plan del MVP (ideas por hacer)

Cada idea es unitaria: un resultado verificable, una capa, una feature, ~300 líneas. Las filas
sin enlace aún no tienen `idea.md`; se crea con `/idea` cuando les toca. El orden numérico
respeta las dependencias. ⛔ marca una decisión abierta que la bloquea.

### A. Catálogo y cartolas

| # | Idea | Estado |
|---|---|---|
| FOKKO-0008-SUBS | Enlace de cancelación en el catálogo | todo |
| FOKKO-0009-SUBS | `extract`: texto de un PDF | todo |
| FOKKO-0010-SUBS | Detectar el banco de un archivo | todo |
| FOKKO-0011-SUBS | Parser BCI: cobros del estado nacional | todo |
| FOKKO-0012-SUBS | Parser BCI: cobros del estado internacional (origen y US$) | todo |
| FOKKO-0013-SUBS | Parser BCI: fecha de facturación | todo |
| FOKKO-0014-SUBS | Parser BCI: par de traspasos | todo |
| FOKKO-0015-SUBS | Derivar el tipo de cambio de los traspasos | todo |
| FOKKO-0016-SUBS | Último tipo de cambio o total incompleto | todo |
| FOKKO-0017-SUBS | Decidir si un cobro internacional es en pesos | todo |
| FOKKO-0018-SUBS | Convertir un cobro en US$ a CLP | todo |
| FOKKO-0019-SUBS | Caso de uso `ImportStatement` | todo |

### B. Detección

| # | Idea | Estado |
|---|---|---|
| FOKKO-0020-SUBS | `UserID` en `Charge` y los puertos de `statements` | todo |
| FOKKO-0021-SUBS | `UserID` en los puertos de `subscriptions` | todo |
| FOKKO-0022-SUBS | Deshacer una decisión | todo |
| FOKKO-0023-SUBS | Próximo cobro por línea | todo |
| FOKKO-0024-SUBS | Cobros ya cobrados del mes | todo |
| FOKKO-0025-SUBS | Cobros restantes del mes | todo |
| FOKKO-0026-SUBS | Detectar un alza | todo |
| FOKKO-0027-SUBS | Impacto mensual y anual del alza | todo |
| FOKKO-0028-SUBS | Descartar el aviso de alza | todo |
| FOKKO-0029-SUBS | Total del mes en CLP | todo |
| FOKKO-0030-SUBS | Monto original junto al equivalente ("≈") | todo |
| FOKKO-0031-SUBS | Crear suscripción manual | todo |
| FOKKO-0032-SUBS | Editar o quitar una manual | todo |
| FOKKO-0033-SUBS | Override de monto y día | todo |
| FOKKO-0034-SUBS | "Dejar de seguir" una línea | todo |
| FOKKO-0035-SUBS | Marcar "por cancelar" | todo |
| FOKKO-0036-SUBS | `Summary`: total y reparto por servicio | todo |
| FOKKO-0037-SUBS | `Summary`: acciones pendientes | todo |
| FOKKO-0038-SUBS | `Summary`: próximos cobros | todo |
| FOKKO-0039-SUBS | `Summary`: detalle de una línea | todo |

### C. Acceso

| # | Idea | Estado |
|---|---|---|
| FOKKO-0040-SUBS | Registrar usuario | todo |
| FOKKO-0041-SUBS | Pedir código al entrar (login) | todo |
| FOKKO-0042-SUBS | Generar el código | todo |
| FOKKO-0043-SUBS | Verificar el código (vigencia y 5 intentos) | todo |
| FOKKO-0044-SUBS | Límite de 3 envíos por hora | todo |
| FOKKO-0045-SUBS | Puerto de envío con doble de consola | todo |
| FOKKO-0046-SUBS | Crear sesión de 30 días | todo |
| FOKKO-0047-SUBS | Validar sesión y cerrarla | todo |
| FOKKO-0048-SUBS | Envío del código con Resend (`net/http`, timeout, llave por entorno) | todo |

### D. Persistencia

| # | Idea | Estado |
|---|---|---|
| FOKKO-0049-SUBS | `sqlc.yaml` y migración base | todo |
| FOKKO-0050-SUBS | Repo de usuarios | todo |
| FOKKO-0051-SUBS | Repo de códigos | todo |
| FOKKO-0052-SUBS | Repo de sesiones | todo |
| FOKKO-0053-SUBS | Repo de cargos (idempotente) | todo |
| FOKKO-0054-SUBS | Repo de decisiones | todo |
| FOKKO-0055-SUBS | Repo de alias | todo |
| FOKKO-0056-SUBS | Repo de overrides | todo |
| FOKKO-0057-SUBS | Repo de manuales | todo |
| FOKKO-0058-SUBS | Repo de líneas ocultas y cancelaciones | todo |
| FOKKO-0059-SUBS | Repo de decisiones de alza | todo |
| FOKKO-0060-SUBS | Repo del tipo de cambio | todo |
| FOKKO-0061-SUBS | Limpiar códigos y sesiones vencidos | todo |

### E. HTTP

| # | Idea | Estado |
|---|---|---|
| FOKKO-0062-SUBS | Base chi y timeouts | todo |
| FOKKO-0063-SUBS | Errores con llaves i18n | todo |
| FOKKO-0064-SUBS | Middleware de sesión | todo |
| FOKKO-0065-SUBS | Cookie (flags), CSRF y CORS | todo |
| FOKKO-0066-SUBS | `register` y `login` | todo |
| FOKKO-0067-SUBS | `code/verify` y `code/resend` | todo |
| FOKKO-0068-SUBS | `logout` y `GET /api/me` | todo |
| FOKKO-0069-SUBS | Límites de la subida (tamaño, PDF, timeout) | todo |
| FOKKO-0070-SUBS | `POST /api/statements` | todo |
| FOKKO-0071-SUBS | `GET /api/banks` | todo |
| FOKKO-0072-SUBS | `GET /api/summary` | todo |
| FOKKO-0073-SUBS | Decisiones (`POST` y `DELETE`) | todo |
| FOKKO-0074-SUBS | Alias y overrides | todo |
| FOKKO-0075-SUBS | Manuales | todo |
| FOKKO-0076-SUBS | Descartar aviso de alza | todo |
| FOKKO-0077-SUBS | `cancel-guide` y cancelaciones | todo |
| FOKKO-0078-SUBS | `DELETE /api/subscriptions/line` | todo |
| FOKKO-0079-SUBS | Exportar copia (CSV/JSON) | todo |
| FOKKO-0080-SUBS | `DELETE /api/me` | todo |
| FOKKO-0081-SUBS | `GET /healthz` | todo |
| FOKKO-0082-SUBS | Armado en `main.go` | todo |

### F. Front

| # | Idea | Estado |
|---|---|---|
| FOKKO-0083-SUBS | Vite, TS estricto y Biome | todo |
| FOKKO-0084-SUBS | Tailwind y tokens de tema | todo |
| FOKKO-0085-SUBS | Router, Query y `api/client` | todo |
| FOKKO-0086-SUBS | i18n y test de contrato | todo |
| FOKKO-0087-SUBS | `Button`, `Input`, `FormField`, `Checkbox` | todo |
| FOKKO-0088-SUBS | `Card`, `Badge`, `ListRow`, `ServiceAvatar` | todo |
| FOKKO-0089-SUBS | `SegmentedRing` y `AnimatedNumber` | todo |
| FOKKO-0090-SUBS | `Toast`, `Skeleton`, `SegmentedControl`, `EmptyState` | todo |
| FOKKO-0091-SUBS | `CodeInput` | todo |
| FOKKO-0092-SUBS | `AppShell` móvil | todo |
| FOKKO-0093-SUBS | `SideNav` de escritorio | todo |
| FOKKO-0094-SUBS | Pantalla Registro | todo |
| FOKKO-0095-SUBS | Pantalla Login | todo |
| FOKKO-0096-SUBS | Pantalla Código | todo |
| FOKKO-0097-SUBS | Pantalla Bienvenida | todo |
| FOKKO-0098-SUBS | Protección de rutas | todo |
| FOKKO-0099-SUBS | Subir cartola | todo |
| FOKKO-0100-SUBS | Cargando y Análisis | todo |
| FOKKO-0101-SUBS | Error de lectura | todo |
| FOKKO-0102-SUBS | Revisión (con selector de cadencia) | todo |
| FOKKO-0103-SUBS | Resumen | todo |
| FOKKO-0104-SUBS | `CycleTip` con el texto corregido | todo |
| FOKKO-0105-SUBS | Subs | todo |
| FOKKO-0106-SUBS | Detalle (con y sin alza) | todo |
| FOKKO-0107-SUBS | Editar y Agregar | todo |
| FOKKO-0108-SUBS | Alerta de alza | todo |
| FOKKO-0109-SUBS | Cómo cancelar | todo |
| FOKKO-0110-SUBS | Perfil y Apariencia | todo |
| FOKKO-0111-SUBS | Eliminar todo | todo |

### G. Cierre

| # | Idea | Estado |
|---|---|---|
| FOKKO-0112-SUBS | Smoke test de punta a punta (subir cartola hasta el Resumen) | todo |
| FOKKO-0113-SUBS | Textos legales y prueba de que se descarta el texto de la cartola ⛔ decisión 12 | todo |
| FOKKO-0114-SUBS | Build único (front servido por el binario) | todo |
| FOKKO-0115-SUBS | Dockerfile multi-stage | todo |
| FOKKO-0116-SUBS | `docker-compose` local | todo |
| FOKKO-0117-SUBS | README | todo |

### Post-MVP (sin número aún)

Lectura de correo por IMAP (cuentas, sincronización, avisos BCI, conciliación), `later` de la
alerta de alza, recordatorio "2 días antes", CSV, OAuth de Gmail, login con Google y Apple,
calendario, notificaciones, recuperar acceso y uso sin conexión.
