# FOKKO-0050-SUBS · Repo de usuarios

## Problema
`auth` necesita guardar usuarios.

## Depende de
- FOKKO-0040-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `users` y adaptador sqlite de `UserRepository`.

## Fuera
- Borrar un usuario y todo lo suyo → FOKKO-0080-SUBS.

## Features afectadas
- `internal/auth` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Tabla `users`](TASK-001-USERS-MIGRATION.md)
- [ ] [TASK-002 Crear usuarios](TASK-002-USERS-REPO-CREATE.md)
- [ ] [TASK-003 Buscar usuarios](TASK-003-USERS-REPO-FIND.md)
