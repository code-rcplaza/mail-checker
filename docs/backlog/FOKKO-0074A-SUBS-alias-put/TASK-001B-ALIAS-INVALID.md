# FOKKO-0074A-SUBS/TASK-001B · Rechazar alias largo o línea inexistente
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dado un alias de más de 40 caracteres, cuando se envía, entonces responde 400 con `subscriptions.alias_too_long`.
- Dada una línea inexistente, cuando se envía, entonces responde 404 con `subscriptions.line_not_found`.

## Fuera
- Nada.

## Notas
- Ninguna.
