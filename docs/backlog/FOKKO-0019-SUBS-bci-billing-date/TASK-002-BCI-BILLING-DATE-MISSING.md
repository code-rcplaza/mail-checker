# FOKKO-0019-SUBS/TASK-002 · Sin fecha en el archivo
Feature: `internal/statements/adapters/bci`

Como usuario quiero que la app me pida la fecha solo cuando no la encuentra, para no hacerlo a ciegas.

## Criterios
- Dada una cartola sin la etiqueta de fecha, cuando se llama `BillingDate`, entonces devuelve `false`.
- Dada una fecha ilegible, cuando se llama `BillingDate`, entonces devuelve `false` y no hace panic.

## Fuera
- Nada.

## Notas
- Ninguna.
