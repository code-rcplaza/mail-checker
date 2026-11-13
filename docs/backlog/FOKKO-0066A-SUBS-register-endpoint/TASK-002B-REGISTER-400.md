# FOKKO-0066A-SUBS/TASK-002B · Responder 400 por datos inválidos
Feature: `internal/auth/adapters/http`



## Criterios
- Dados campos faltantes o inválidos, cuando se envía, entonces responde 400 con `common.validation_failed` y `fields`.
- Dado un cuerpo con un campo desconocido, cuando se envía, entonces responde 400 con `common.invalid_json`.

## Fuera
- Nada.

## Notas
- Ninguna.
