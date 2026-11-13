# FOKKO-0016-SUBS/TASK-002 · Sin tipo de cambio
Feature: `internal/subscriptions/domain`

Como usuario quiero que la app me avise si no puede convertir mis dólares, para no confiar en un total incompleto.

## Criterios
- Dada una lista vacía, cuando se elige, entonces devuelve que no hay tipo de cambio.
- Dada esa respuesta, cuando se arma el total, entonces el resultado queda marcado como incompleto (lo usa FOKKO-0029-SUBS).

## Fuera
- Nada.

## Notas
- Ninguna.
