# FOKKO-0033-SUBS/TASK-002 · Guardar o quitar un override
Feature: `internal/subscriptions/usecases`

Como usuario quiero corregir el monto o el día de una línea, para que mi resumen refleje lo que pago.

## Criterios
- Dada una línea existente y un override válido, cuando se guarda, entonces se guarda por clave de línea.
- Dado un override vacío, cuando se guarda, entonces se borra el que hubiera.
- Dada una línea inexistente o de otro usuario, cuando se guarda, entonces devuelve `ErrNotFound` con `subscriptions.line_not_found`.

## Fuera
- Nada.

## Notas
- Ninguna.
