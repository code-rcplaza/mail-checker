# FOKKO-0074B-SUBS · `PUT /api/subscriptions/overrides`

## Problema
Editar permite nombrar una línea y reemplazar su monto o día.

Esta idea cubre solo: `PUT /api/subscriptions/overrides`. El resto del trabajo original está en las ideas hermanas `FOKKO-0074*-SUBS`.

## Depende de
- FOKKO-0033-SUBS y FOKKO-0064-SUBS.

## Alcance
- `PUT /api/subscriptions/overrides`

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Guardar o borrar el override](TASK-001A-OVERRIDE-SET.md)
- [ ] [TASK-001B Rechazar override inválido o línea inexistente](TASK-001B-OVERRIDE-INVALID.md)
