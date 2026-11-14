# FOKKO-0053A-SUBS · Repo de cartolas: crear con sus cargos

## Problema
Una cartola y sus cargos deben guardarse juntos o ninguno, sin duplicar al reimportar.

## Depende de
- FOKKO-0020A-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `statements` y `charges` (compartida) y creación atómica de una cartola con sus cargos.

## Fuera
- Listar y borrar → FOKKO-0053B-SUBS.
- Leer cargos desde `subscriptions` → FOKKO-0053C-SUBS.

## Features afectadas
- `internal/statements` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Dentro de `sqlite.WithTx` solo se usa la transacción, nunca el repo ni `db` (la base tiene una sola conexión y se bloquea). Los tests con `-timeout 60s` lo detectan.

## Tareas
- [x] [TASK-001A Tabla `statements`](TASK-001A-MIGRATION-STATEMENTS.md)
- [x] [TASK-001B Tabla `charges`](TASK-001B-MIGRATION-CHARGES.md)
- [x] [TASK-002A Crear una cartola con sus cargos](TASK-002A-CREATE-WITH-CHARGES.md)
- [x] [TASK-002B Rechazar una cartola repetida](TASK-002B-CREATE-CONFLICT.md)
