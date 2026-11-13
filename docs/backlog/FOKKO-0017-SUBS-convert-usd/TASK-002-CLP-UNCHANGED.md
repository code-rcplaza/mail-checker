# FOKKO-0017-SUBS/TASK-002 · Un monto en pesos no se toca
Feature: `internal/subscriptions/domain`

Como usuario quiero que un cobro en pesos salga igual, para no ver diferencias inventadas.

## Criterios
- Dado un monto en CLP, cuando se convierte, entonces devuelve el mismo monto y no es aproximado.
- Dada una moneda distinta de CLP y USD, cuando se convierte, entonces devuelve `ErrInvalid`.

## Fuera
- Nada.

## Notas
- Ninguna.
