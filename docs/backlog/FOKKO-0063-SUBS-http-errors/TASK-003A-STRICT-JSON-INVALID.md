# FOKKO-0063-SUBS/TASK-003A · Rechazar JSON inválido o vacío
Feature: `internal/shared/httpserver`



## Criterios
- Dado un cuerpo que no es JSON, cuando se decodifica, entonces devuelve `ErrInvalid` con `common.invalid_json`.
- Dado un cuerpo vacío donde se esperaba JSON, cuando se decodifica, entonces devuelve `common.invalid_json`.

## Fuera
- Nada.

## Notas
- Ninguna.
