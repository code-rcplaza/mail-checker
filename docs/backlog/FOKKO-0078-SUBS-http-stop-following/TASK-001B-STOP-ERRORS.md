# FOKKO-0078-SUBS/TASK-001B · Rechazar clave ausente o línea inexistente
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dada una línea inexistente, cuando se envía, entonces responde 404 con `subscriptions.line_not_found`.
- Dada una clave sin `lineKey`, cuando se envía, entonces responde 400 con `fields`.

## Fuera
- Nada.

## Notas
- Ninguna.
