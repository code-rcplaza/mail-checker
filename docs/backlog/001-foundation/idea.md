# 001 · Armazón del backend (fase 1)

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

## Historias
- [x] [01 Módulo, errores de dominio y reloj](01-module-errors-clock.md)
- [x] [02 Configuración con fail fast](02-config.md)
- [ ] [03 Abrir SQLite con PRAGMAs](03-sqlite-open.md)
- [ ] [04 Migraciones embebidas](04-migrations.md)
- [ ] [05 Transacciones con `WithTx`](05-with-tx.md)
- [ ] [06 Router, health y JSON](06-router-health-json.md)
- [ ] [07 Mapeo de errores a HTTP](07-error-mapping.md)
- [ ] [08 `main.go`, `slog` y apagado ordenado](08-main-shutdown.md)
