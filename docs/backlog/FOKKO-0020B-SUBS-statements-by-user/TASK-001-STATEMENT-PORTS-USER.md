# FOKKO-0020B-SUBS/TASK-001 · Puertos de `statements` por usuario
Feature: `internal/statements/ports`

Como usuario quiero ver solo mis cartolas, para que mis datos sean privados.

## Criterios
- Dado `StatementRepository`, cuando se llama `Create`, `List` o `Delete`, entonces recibe el `userID`.
- Dados dos usuarios con cartolas, cuando uno lista, entonces el doble en memoria devuelve solo las suyas.
- Dado un `Delete` con la cartola de otro usuario, cuando se ejecuta, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
