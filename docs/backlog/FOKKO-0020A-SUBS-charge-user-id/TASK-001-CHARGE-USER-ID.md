# FOKKO-0020A-SUBS/TASK-001 · `UserID` en `Charge`
Feature: `internal/shared/domain`

Como usuario quiero que cada cargo tenga dueño, para que nunca se mezcle con el de otra persona.

## Criterios
- Dado un cargo con `UserID` mayor que cero, cuando se valida con `NewCharge`, entonces es válido.
- Dado un cargo con `UserID` cero o negativo, cuando se valida, entonces devuelve `ErrInvalid`.
- Dado `ChargeOrigin`, cuando se convierte un `ParsedCharge`, entonces el cargo lleva el `UserID` del origen.

## Fuera
- Nada.

## Notas
- Ninguna.
