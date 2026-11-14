# FOKKO-0054-SUBS · Repo de decisiones

## Problema
Las decisiones del usuario sobre las posibles deben persistir por clave de línea.

## Depende de
- FOKKO-0021-SUBS, FOKKO-0022-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `subscription_decisions` y adaptador (`List`, `Save`, `Delete`).

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Tabla `subscription_decisions`](TASK-001-DECISIONS-MIGRATION.md)
- [x] [TASK-002A Guardar una decisión](TASK-002A-DECISIONS-SAVE.md)
- [x] [TASK-002B Listar mis decisiones](TASK-002B-DECISIONS-LIST.md)
- [x] [TASK-002C Borrar una decisión](TASK-002C-DECISIONS-DELETE.md)
