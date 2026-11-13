# FOKKO-0043-SUBS/TASK-002 · Código equivocado
Feature: `internal/auth/usecases`

Como usuario quiero que se cuenten mis intentos fallidos, para que nadie adivine mi código.

## Criterios
- Dado un código equivocado, cuando se verifica, entonces devuelve `ErrInvalid` con `auth.code_invalid` y suma un intento.
- Dado el quinto intento fallido, cuando se verifica, entonces el código queda invalidado y devuelve `auth.code_attempts_exceeded`.
- Dado un código invalidado, cuando se verifica el correcto, entonces sigue rechazándolo.

## Fuera
- Nada.

## Notas
- Ninguna.
