# FOKKO-0063-SUBS · Errores con llaves i18n

## Problema
El back no devuelve textos: devuelve llaves i18n. Hoy se mapean 400, 404, 409 y 500; los nuevos errores de sesión (`ErrUnauthorized`) y de límite (`ErrTooMany`) necesitan su status, y los errores de campo necesitan `fields`.

## Depende de
- FOKKO-0044-SUBS y FOKKO-0047A-SUBS y FOKKO-0047B-SUBS.

## Alcance
- Mapear `ErrUnauthorized` a 401 y `ErrTooMany` a 429.
- Respuestas de validación por campo (`common.validation_failed` con `fields`).
- Decodificación estricta del JSON de entrada.

## Fuera
- Que el front tenga todas las llaves → FOKKO-0086A-SUBS y FOKKO-0086B-SUBS.

## Features afectadas
- `internal/shared` (`httpserver`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Responder 401 por falta de sesión](TASK-001A-MAP-401.md)
- [x] [TASK-001B Responder 429 por límite](TASK-001B-MAP-429.md)
- [x] [TASK-001C Responder 500 sin detalle](TASK-001C-MAP-500.md)
- [x] [TASK-002 Errores por campo](TASK-002-FIELD-ERRORS.md)
- [x] [TASK-003A Rechazar JSON inválido o vacío](TASK-003A-STRICT-JSON-INVALID.md)
- [x] [TASK-003B Rechazar campos o valores de más](TASK-003B-STRICT-JSON-EXTRA.md)
