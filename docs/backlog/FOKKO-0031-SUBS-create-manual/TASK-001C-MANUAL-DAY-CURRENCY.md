# FOKKO-0031-SUBS/TASK-001C · Rechazar día o moneda inválidos
Feature: `internal/subscriptions/domain`



## Criterios
- Dado un día 0 o 32, cuando se construye, entonces devuelve `ErrInvalid`.
- Dada una moneda distinta de CLP y USD, cuando se construye, entonces devuelve `ErrInvalid`.

## Fuera
- Nada.

## Notas
- Ninguna.
