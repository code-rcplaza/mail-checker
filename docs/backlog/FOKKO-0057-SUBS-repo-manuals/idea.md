# FOKKO-0057-SUBS · Repo de manuales

## Problema
Las suscripciones manuales deben persistir y exponerse con un `public_id`.

## Depende de
- FOKKO-0031-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `manual_subscriptions` y adaptador (`Create`, `Update`, `Delete`, `List`).

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Tabla `manual_subscriptions`](TASK-001-MANUALS-MIGRATION.md)
- [ ] [TASK-002A Crear una manual](TASK-002A-MANUALS-CREATE.md)
- [ ] [TASK-002B Listar mis manuales](TASK-002B-MANUALS-LIST.md)
- [ ] [TASK-003A Editar una manual](TASK-003A-MANUALS-UPDATE.md)
- [ ] [TASK-003B Quitar una manual](TASK-003B-MANUALS-DELETE.md)
