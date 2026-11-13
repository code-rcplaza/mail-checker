# FOKKO-0051A-SUBS · Repo de códigos de acceso

## Problema
Los códigos de acceso deben persistir para sobrevivir a un reinicio. Un solo código vigente por usuario, con su vencimiento e intentos.

## Depende de
- FOKKO-0042-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `login_codes` y adaptador sqlite del puerto de códigos.

## Fuera
- Registro de envíos → FOKKO-0051B-SUBS.
- Limpiar vencidos → FOKKO-0061-SUBS.

## Features afectadas
- `internal/auth` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Tabla `login_codes`](TASK-001-CODES-MIGRATION.md)
- [ ] [TASK-002A Guardar un código](TASK-002A-CODES-SAVE.md)
- [ ] [TASK-002B Sumar un intento](TASK-002B-CODES-ATTEMPTS.md)
- [ ] [TASK-002C Consumir un código](TASK-002C-CODES-CONSUME.md)
