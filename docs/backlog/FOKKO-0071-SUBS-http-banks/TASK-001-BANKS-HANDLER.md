# FOKKO-0071-SUBS/TASK-001 · `GET /api/banks`
Feature: `internal/statements/adapters/http`

Como usuario quiero elegir mi banco si la app no lo detecta, para que lea mi cartola.

## Criterios
- Dado el registro con BCI, cuando se pide, entonces responde 200 con `[{id:"bci"}]`.
- Dado un registro vacío, cuando se pide, entonces responde 200 con `[]`.
- Dada una petición sin sesión, cuando se pide, entonces responde 401.

## Fuera
- Nada.

## Notas
- Ninguna.
