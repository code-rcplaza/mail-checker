# FOKKO-0047A-SUBS/TASK-001 · Error de falta de sesión
Feature: `internal/shared/domain`

Como usuario quiero saber que debo entrar de nuevo, para no ver un error genérico.

## Criterios
- Dada una llave, cuando se crea el error, entonces envuelve `ErrUnauthorized` y lleva la llave (`common.unauthorized`).
- Dado ese error, cuando se compara con `errors.Is`, entonces no coincide con `ErrNotFound` ni `ErrInvalid`.

## Fuera
- Nada.

## Notas
- Ninguna.
