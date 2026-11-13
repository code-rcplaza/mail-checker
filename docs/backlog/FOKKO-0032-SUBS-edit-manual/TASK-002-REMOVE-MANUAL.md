# FOKKO-0032-SUBS/TASK-002 · Quitar una manual
Feature: `internal/subscriptions/usecases`

Como usuario quiero quitar una manual que ya no pago, para que no sume al total.

## Criterios
- Dada una manual existente, cuando se quita, entonces ya no aparece en el resumen.
- Dado un `PublicID` inexistente o de otro usuario, cuando se quita, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
