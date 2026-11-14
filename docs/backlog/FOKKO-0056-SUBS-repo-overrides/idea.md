# FOKKO-0056-SUBS · Repo de overrides

## Problema
Los overrides de monto y día de una línea deben persistir.

## Depende de
- FOKKO-0033-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `line_overrides` y adaptador (`List`, `Save`, `Delete`).

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Tabla `line_overrides`](TASK-001-OVERRIDES-MIGRATION.md)
- [x] [TASK-002 Adaptador de overrides](TASK-002-OVERRIDES-REPO.md)
