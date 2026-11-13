# FOKKO-0047A-SUBS/TASK-002 · Validar una sesión
Feature: `internal/auth/usecases`

Como usuario quiero seguir dentro mientras mi sesión valga, para no entrar a cada rato.

## Criterios
- Dado un token vigente, cuando se valida, entonces devuelve el `userID`.
- Dado un token vencido, cuando se valida, entonces devuelve `ErrUnauthorized`.
- Dado un token desconocido, cuando se valida, entonces devuelve `ErrUnauthorized` sin distinguir el motivo.

## Fuera
- Nada.

## Notas
- Ninguna.
