# FOKKO-0078-SUBS · `DELETE /api/subscriptions/line`

## Problema
"Dejar de seguir" oculta una línea desde Detalle.

## Depende de
- FOKKO-0034-SUBS y FOKKO-0064-SUBS.

## Alcance
- `DELETE /api/subscriptions/line?lineKey=`.

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Dejar de seguir una línea](TASK-001A-STOP-OK.md)
- [ ] [TASK-001B Rechazar clave ausente o línea inexistente](TASK-001B-STOP-ERRORS.md)
