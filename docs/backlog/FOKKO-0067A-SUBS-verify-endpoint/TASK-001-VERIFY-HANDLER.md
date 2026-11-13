# FOKKO-0067A-SUBS/TASK-001 · `POST /api/auth/code/verify`
Feature: `internal/auth/adapters/http`

Como usuario quiero entrar al poner mi código, para ver mis suscripciones.

## Criterios
- Dado `{email, code}` correcto, cuando se envía, entonces responde 204 y fija la cookie de sesión.
- Dado un código incorrecto, cuando se envía, entonces responde 400 con `auth.code_invalid` y no fija cookie.
- Dado un código vencido o agotado, cuando se envía, entonces responde 400 con `auth.code_expired` o `auth.code_attempts_exceeded`.

## Fuera
- Nada.

## Notas
- Ninguna.
