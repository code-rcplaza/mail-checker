# FOKKO-0052-SUBS · Repo de sesiones

## Problema
Las sesiones deben persistir entre reinicios, guardadas por el hash del token.

## Depende de
- FOKKO-0046-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `sessions` y adaptador sqlite.

## Fuera
- Nada.

## Features afectadas
- `internal/auth` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Tabla `sessions`](TASK-001-SESSIONS-MIGRATION.md)
- [ ] [TASK-002A Buscar una sesión por su hash](TASK-002A-SESSIONS-FIND.md)
- [ ] [TASK-002B Borrar una sesión](TASK-002B-SESSIONS-DELETE.md)
