# FOKKO-0063-SUBS/TASK-003B · Rechazar campos o valores de más
Feature: `internal/shared/httpserver`



## Criterios
- Dado un cuerpo con un campo desconocido, cuando se decodifica, entonces devuelve `common.invalid_json` en vez de ignorarlo.
- Dado un cuerpo con más de un valor JSON, cuando se decodifica, entonces devuelve `common.invalid_json`.

## Fuera
- Nada.

## Notas
- Ninguna.
