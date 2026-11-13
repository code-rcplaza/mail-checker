# FOKKO-0003-SUBS · Núcleo de sincronización de correos (fase 2, parte 2)

## Problema
La detección ya tiene de dónde leer cobros, pero nada los trae: falta el flujo que
recorre las cuentas, busca los avisos de cada banco, los convierte en cobros y deja
registro de lo que no se pudo leer. Esta idea implementa ese flujo con TDD y dobles
en memoria, sin IMAP ni base.

## Alcance
- `shared/domain`: `ParsedCharge` (lo que un parser extrae de un aviso).
- `emailsync/domain`:
  - `RawEmail`, `MailQuery` (filtro del banco + desde cuándo), `Credential`,
    `SyncResult`, `UnparsedEmail`.
  - Conversión `ParsedCharge` → `Charge`: banco del parser, fuente `email`, external id
    = Message-ID, cuenta; si el aviso no trae fecha, se usa la fecha del correo.
  - Llaves: `emailsync.mail_auth_failed` (credenciales rechazadas),
    `emailsync.mail_read_failed` (otro error de lectura), `emailsync.sync_in_progress`,
    `emailsync.months_out_of_range`.
- `emailsync/ports`: `AccountDirectory`, `MailReader`, `BankEmailParser`,
  `ChargeWriter`, `UnparsedRepository` (dobles en memoria).
- `emailsync/usecases`:
  - `SyncAll(months)`: para cada cuenta y cada banco registrado, busca con el filtro del
    banco desde hoy − months, parsea y guarda.
    - El mismo aviso en dos cuentas se guarda una vez (Message-ID); `New` cuenta solo
      los insertados.
    - Aviso que el parser no reconoce → `UnparsedRepository.Save`; si después se
      reconoce → `Forget`.
    - Un error en una cuenta queda en su `SyncResult` (`ErrorKey`) y no detiene a las
      demás; se registra con `RecordSync`.
    - Una sola sincronización a la vez: otra en curso → `ErrConflict` con
      `emailsync.sync_in_progress`.
    - `months` fuera de 1–24 → `ErrInvalid` (mismo rango que el resumen).
  - `ListUnparsed(limit)`.

## Fuera
- Parser real de BCI → fase 3.
- IMAP (`MailReader` real) y detección real de credenciales rechazadas → fase 5.
- Repositorios SQLite (`ChargeWriter`, `UnparsedRepository`) → fase 4.
- Puente a `accounts` (`AccountDirectory` real) → con la idea de `accounts`.
- Endpoints `/api/sync` y `/api/unparsed` → fase 5.

## Features afectadas
- `internal/shared/domain`
- `internal/emailsync`

## Dudas abiertas
- Ninguna. Depende de FOKKO-0002-SUBS/TASK-001 (`Charge`, `NewCharge`). Agregados al aprobar las
  historias: `ErrMailAuthRejected` y `emailsync.limit_out_of_range`.
- Corrección menor en `ARCHITECTURE.md`: el ejemplo de `SyncResult.ErrorKey`
  dice `accounts.imap_auth_failed`; la llave es de `emailsync`
  (`emailsync.mail_auth_failed`).

## Tareas
- [x] [TASK-001 `ParsedCharge` y su conversión a `Charge`](TASK-001-PARSED-CHARGE.md)
- [x] [TASK-002 Tipos de `emailsync`](TASK-002-EMAILSYNC-TYPES.md)
- [x] [TASK-003 `SyncAll`: camino feliz](TASK-003-SYNC-ALL-HAPPY-PATH.md)
- [x] [TASK-004 `SyncAll`: duplicados y no reconocidos](TASK-004-SYNC-ALL-DUPLICATES-UNPARSED.md)
- [x] [TASK-005 `SyncAll`: errores por cuenta](TASK-005-SYNC-ALL-ACCOUNT-ERRORS.md)
- [x] [TASK-006 `SyncAll`: validaciones](TASK-006-SYNC-ALL-GUARDS.md)
- [x] [TASK-007 `ListUnparsed`](TASK-007-LIST-UNPARSED.md)
