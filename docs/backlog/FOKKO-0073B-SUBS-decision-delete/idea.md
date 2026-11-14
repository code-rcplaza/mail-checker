# FOKKO-0073B-SUBS · `DELETE /api/subscriptions/decisions`

## Problema
Revisión confirma (con cadencia) o descarta una posible, y puede deshacerlo.

Esta idea cubre solo: `DELETE /api/subscriptions/decisions`. El resto del trabajo original está en las ideas hermanas `FOKKO-0073*-SUBS`.

## Depende de
- FOKKO-0022-SUBS y FOKKO-0064-SUBS.

## Alcance
- `DELETE /api/subscriptions/decisions`

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Deshacer una decisión](TASK-001A-DECISION-UNDO.md)
- [x] [TASK-001B Validar la clave de línea](TASK-001B-DECISION-KEY.md)
