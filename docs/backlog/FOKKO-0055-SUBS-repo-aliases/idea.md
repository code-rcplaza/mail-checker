# FOKKO-0055-SUBS · Repo de alias

## Problema
Los nombres que el usuario da a sus líneas deben persistir.

## Depende de
- FOKKO-0021-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `line_aliases` y adaptador (`List`, `Save`, `Delete`).

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Tabla `line_aliases`](TASK-001-ALIASES-MIGRATION.md)
- [x] [TASK-002A Guardar un alias](TASK-002A-ALIASES-SAVE.md)
- [x] [TASK-002B Listar mis alias](TASK-002B-ALIASES-LIST.md)
- [x] [TASK-002C Borrar un alias](TASK-002C-ALIASES-DELETE.md)
