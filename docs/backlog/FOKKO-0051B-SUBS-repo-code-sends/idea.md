# FOKKO-0051B-SUBS · Repo del registro de envíos de códigos

## Problema
El límite de 3 envíos por hora por correo necesita un registro persistente de envíos.

## Depende de
- FOKKO-0044-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `code_sends` y adaptador para registrar y contar envíos.

## Fuera
- Los códigos → FOKKO-0051A-SUBS.

## Features afectadas
- `internal/auth` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Tabla `code_sends`](TASK-001-SENDS-MIGRATION.md)
- [ ] [TASK-002A Registrar un envío](TASK-002A-SENDS-RECORD.md)
- [ ] [TASK-002B Contar los envíos de la última hora](TASK-002B-SENDS-COUNT.md)
