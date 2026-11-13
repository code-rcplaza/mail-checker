# FOKKO-0031-SUBS/TASK-001B · Rechazar nombre o monto inválidos
Feature: `internal/subscriptions/domain`



## Criterios
- Dado un nombre vacío, cuando se construye, entonces devuelve `ErrInvalid` con su llave.
- Dado un monto cero o negativo, cuando se construye, entonces devuelve `ErrInvalid`.

## Fuera
- Nada.

## Notas
- Ninguna.
