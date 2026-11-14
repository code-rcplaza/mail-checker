# FOKKO-0085B-SUBS · Cimientos de TanStack (Query, Router)

## Problema
Todos los datos del servidor viven en TanStack Query y se leen con un cliente que entiende el formato de error del back (llaves i18n, parámetros, campos).

Esta idea cubre solo: Cimientos de TanStack (Query, Router). El resto del trabajo original está en las ideas hermanas `FOKKO-0085*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS.

## Alcance
- Query y router

## Fuera
- Las llaves de texto → FOKKO-0086A-SUBS y FOKKO-0086B-SUBS.

## Features afectadas
- `frontend` (`lib/api`, `routes`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Reintentos de Query](TASK-001A-QUERY-DEFAULTS.md)
- [x] [TASK-001B Pantalla 404 del router](TASK-001B-ROUTER-404.md)
- [x] [TASK-001C Sin estado global duplicado](TASK-001C-NO-GLOBAL-STATE.md)
