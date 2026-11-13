# FOKKO-0066B-SUBS/TASK-001B · Rechazar correo inválido o desconocido
Feature: `internal/auth/adapters/http`



## Criterios
- Dado un correo con formato inválido, cuando se envía, entonces responde 400 con `fields`.
- Dado un correo no registrado, cuando se envía, entonces responde 404 con `auth.user_not_found`.

## Fuera
- Nada.

## Notas
- Ninguna.
