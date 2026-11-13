# FOKKO-0058B-SUBS · Repo de cancelaciones

## Problema
La marca "por cancelar" debe persistir con la fecha en que se marcó.

## Depende de
- FOKKO-0035-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `cancellations` y su adaptador.

## Fuera
- Líneas ocultas → FOKKO-0058A-SUBS.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Tabla `cancellations`](TASK-001-CANCELLATIONS-MIGRATION.md)
- [ ] [TASK-002 Adaptador de cancelaciones](TASK-002-CANCELLATIONS-REPO.md)
