# FOKKO-0050-SUBS/TASK-003 · Buscar usuarios
Feature: `internal/auth/adapters/sqlite`

Como usuario quiero que la app me reconozca por mi correo, para entrar.

## Criterios
- Dado un usuario guardado, cuando se busca por correo en otra capitalización, entonces lo encuentra.
- Dado un correo inexistente, cuando se busca, entonces devuelve `ErrNotFound`.
- Dado un usuario guardado, cuando se busca por `ID`, entonces lo encuentra.

## Fuera
- Nada.

## Notas
- Ninguna.
