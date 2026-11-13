# FOKKO-0001-SUBS · Armazón del backend (fase 1)

## Problema
No hay código. Si la primera feature tiene que montar también la base, el router,
la configuración y el logging, queda mezclada con infraestructura. Hace falta un
armazón compartido, probado y reutilizable, para que cada feature después solo
agregue su lógica y sus adaptadores.

## Alcance
- `backend/go.mod` (módulo `github.com/code-rcplaza/mail-checker/backend`, Go 1.27)
  con chi, `modernc.org/sqlite` y goose.
- `shared/sqlite`:
  - `Open`: aplica `journal_mode=WAL`, `foreign_keys=ON`, `busy_timeout=5000`,
    `synchronous=FULL`; un solo escritor.
  - Migraciones goose embebidas, **sin tablas**; si fallan, la app no arranca.
  - `WithTx`: helper para que los repos abran y cierren su transacción.
- `shared/httpserver`:
  - Router chi con `GET /api/health`.
  - Helpers para leer y escribir JSON.
  - Mapeo de errores: `ErrInvalid` → 400, `ErrNotFound` → 404, `ErrConflict` → 409,
    resto → 500; cuerpo `{ "error": { "key", "params" } }` con llave i18n.
  - Apagado ordenado.
- `shared/domain`: solo `ErrInvalid`, `ErrNotFound`, `ErrConflict`, sus constructores con
  llave i18n (`Invalid(key, params)`, `NotFound(key)`, `Conflict(key)`) y las llaves `common.*`.
- `shared/clock`: reloj real.
- Configuración por variables de entorno (puerto, ruta de la base) en
  `cmd/api/config.go`; si es inválida, falla al arrancar.
- Logger `slog`.
- `cmd/api/main.go` mínimo: config → logger → base → router → servidor.
- Tests: `Open` aplica los PRAGMAs; migraciones corren; `WithTx` hace commit y
  rollback; health con `httptest`; mapeo de errores; config inválida falla.

## Fuera
- Tablas, queries y repositorios de cada feature → fase 4.
- `Money`, `Charge` y demás tipos de dominio → llegan con la primera feature (fase 2).
- Front (esqueleto) → fase 6, como idea propia.
- Makefile → fase 10.
- `sqlc.yaml` → fase 4, con el primer repositorio (sin queries no hay qué verificar).
- `google/uuid` → con el primer repo que exponga `public_id`.

## Features afectadas
- `internal/shared/sqlite`
- `internal/shared/httpserver`
- `internal/shared/domain`
- `internal/shared/clock`
- `cmd/api`

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Módulo, errores de dominio y reloj](TASK-001-MODULE-ERRORS-CLOCK.md)
- [x] [TASK-002 Configuración con fail fast](TASK-002-CONFIG.md)
- [x] [TASK-003 Abrir SQLite con PRAGMAs](TASK-003-SQLITE-OPEN.md)
- [x] [TASK-004 Migraciones embebidas](TASK-004-MIGRATIONS.md)
- [x] [TASK-005 Transacciones con `WithTx`](TASK-005-WITH-TX.md)
- [x] [TASK-006 Router, health y JSON](TASK-006-ROUTER-HEALTH-JSON.md)
- [x] [TASK-007 Mapeo de errores a HTTP](TASK-007-ERROR-MAPPING.md)
- [x] [TASK-008 `main.go`, `slog` y apagado ordenado](TASK-008-MAIN-SHUTDOWN.md)
