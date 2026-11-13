# FOKKO-0043-SUBS/TASK-001 · Código correcto
Feature: `internal/auth/usecases`

Como usuario quiero entrar al poner el código correcto, para abrir mi sesión.

## Criterios
- Dado un código vigente y correcto, cuando se verifica, entonces devuelve al usuario y consume el código.
- Dado el mismo código usado una segunda vez, cuando se verifica, entonces devuelve `ErrInvalid` con `auth.code_invalid`.
- Dado un correo en mayúsculas, cuando se verifica, entonces funciona igual.

## Fuera
- Nada.

## Notas
- Ninguna.
