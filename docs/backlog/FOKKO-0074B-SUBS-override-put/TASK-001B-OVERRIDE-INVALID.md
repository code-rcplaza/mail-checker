# FOKKO-0074B-SUBS/TASK-001B · Rechazar override inválido o línea inexistente
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dado un monto cero o negativo, o un día fuera de 1 a 31, cuando se envía, entonces responde 400 con `fields`.
- Dada una línea inexistente, cuando se envía, entonces responde 404.

## Fuera
- Nada.

## Notas
- Ninguna.
