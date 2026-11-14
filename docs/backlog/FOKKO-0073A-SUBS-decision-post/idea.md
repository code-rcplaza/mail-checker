# FOKKO-0073A-SUBS · `POST /api/subscriptions/decisions`

## Problema
Revisión confirma (con cadencia) o descarta una posible, y puede deshacerlo.

Esta idea cubre solo: `POST /api/subscriptions/decisions`. El resto del trabajo original está en las ideas hermanas `FOKKO-0073*-SUBS`.

## Depende de
- FOKKO-0022-SUBS y FOKKO-0064-SUBS.

## Alcance
- `POST /api/subscriptions/decisions`

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Confirmar con cadencia](TASK-001A-DECISION-CONFIRM.md)
- [x] [TASK-001B Descartar una posible](TASK-001B-DECISION-DISMISS.md)
- [x] [TASK-001C Rechazar una clave que no es posible](TASK-001C-DECISION-NOT-POSSIBLE.md)
