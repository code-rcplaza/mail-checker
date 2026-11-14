# FOKKO-0085A-SUBS · Cliente de la API

## Problema
Todos los datos del servidor viven en TanStack Query y se leen con un cliente que entiende el formato de error del back (llaves i18n, parámetros, campos).

Esta idea cubre solo: Cliente de la API. El resto del trabajo original está en las ideas hermanas `FOKKO-0085*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS.

## Alcance
- Cliente de la API
- Errores de la API tipados

## Fuera
- Las llaves de texto → FOKKO-0086A-SUBS y FOKKO-0086B-SUBS.

## Features afectadas
- `frontend` (`lib/api`, `routes`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Leer respuestas correctas](TASK-001A-CLIENT-SUCCESS.md)
- [x] [TASK-001B Leer errores de la API](TASK-001B-CLIENT-ERRORS.md)
- [x] [TASK-001C Enviar cookie con origen propio](TASK-001C-CLIENT-CREDENTIALS.md)
- [x] [TASK-002 Errores de la API tipados](TASK-002-API-ERROR-TYPES.md)
