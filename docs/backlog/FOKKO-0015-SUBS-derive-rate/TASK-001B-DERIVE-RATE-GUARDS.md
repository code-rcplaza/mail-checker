# FOKKO-0015-SUBS/TASK-001B · Rechazar montos inválidos
Feature: `internal/shared/domain`



## Criterios
- Dado US$ 0,00, cuando se deriva, entonces no hay tipo de cambio y no hay división por cero.
- Dado un CLP o US$ negativo, cuando se deriva, entonces devuelve `ErrInvalid`.

## Fuera
- Nada.

## Notas
- Los montos del ejemplo son de prueba, no de una cartola real.
