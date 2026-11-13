# FOKKO-0052-SUBS/TASK-002A · Buscar una sesión por su hash
Feature: `internal/auth/adapters/sqlite`



## Criterios
- Dada una sesión guardada, cuando se busca por hash, entonces devuelve al usuario y el vencimiento.
- Dado un hash desconocido, cuando se busca, entonces devuelve `ErrNotFound`.
- Dada una sesión nueva, cuando se guarda, entonces se puede buscar por su hash.

## Fuera
- Nada.

## Notas
- Ninguna.
