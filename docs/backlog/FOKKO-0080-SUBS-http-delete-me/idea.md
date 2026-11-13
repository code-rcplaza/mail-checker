# FOKKO-0080-SUBS · `DELETE /api/me`

## Problema
"Eliminar todo" borra al usuario y todo lo suyo. Las tablas de datos llevan `user_id` con `ON DELETE CASCADE`, así que basta borrar la fila del usuario en una transacción.

## Depende de
- FOKKO-0050-SUBS a FOKKO-0059-SUBS.

## Alcance
- Borrar al usuario, comprobar la cascada y `DELETE /api/me`.

## Fuera
- La copia previa → FOKKO-0079-SUBS.

## Features afectadas
- `internal/auth` (`ports`, `usecases`, `adapters/sqlite`, `adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Eliminar al usuario](TASK-001-DELETE-USER-USECASE.md)
- [ ] [TASK-002 La cascada borra todo](TASK-002-USERS-REPO-DELETE-CASCADE.md)
- [ ] [TASK-003 `DELETE /api/me`](TASK-003-DELETE-ME-HANDLER.md)
