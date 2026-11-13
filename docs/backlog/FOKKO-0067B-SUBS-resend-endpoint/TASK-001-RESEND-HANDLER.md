# FOKKO-0067B-SUBS/TASK-001 · `POST /api/auth/code/resend`
Feature: `internal/auth/adapters/http`

Como usuario quiero pedir otro código si no me llegó, para poder entrar.

## Criterios
- Dado `{email}` registrado dentro del límite, cuando se envía, entonces responde 204 y el código anterior deja de valer.
- Dado el límite de 3 envíos por hora superado, cuando se envía, entonces responde 429 con `auth.too_many_sends`.

## Fuera
- Nada.

## Notas
- Ninguna.
