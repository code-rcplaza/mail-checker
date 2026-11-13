# FOKKO-0032-SUBS/TASK-001 · Editar una manual
Feature: `internal/subscriptions/usecases`

Como usuario quiero corregir el monto o el día de una manual, para mantenerla al día.

## Criterios
- Dados datos válidos, cuando se edita, entonces se guardan y la manual conserva su `PublicID`.
- Dados datos inválidos, cuando se edita, entonces devuelve `ErrInvalid` y no cambia nada.
- Dado un `PublicID` inexistente o de otro usuario, cuando se edita, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
