# FOKKO-0053A-SUBS/TASK-002B · Rechazar una cartola repetida
Feature: `internal/statements/adapters/sqlite`



## Criterios
- Dada una cartola ya existente (mismo usuario, banco y fecha), cuando se crea, entonces devuelve `ErrConflict`.
- Dado el mismo archivo de otro usuario, cuando se crea, entonces se guarda sin conflicto.

## Fuera
- Nada.

## Notas
- Ninguna.
