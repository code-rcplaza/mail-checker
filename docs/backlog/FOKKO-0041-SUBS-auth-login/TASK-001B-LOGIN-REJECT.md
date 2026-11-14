# FOKKO-0041-SUBS/TASK-001B · Rechazar un correo inválido o desconocido
Feature: `internal/auth/usecases`



## Criterios
- Dado un correo inválido, cuando se pide entrar, entonces devuelve `ErrInvalid` con `auth.email_invalid` sin tocar el repositorio.
- Dado un correo no registrado, cuando se pide entrar, entonces responde sin error, igual que con uno registrado, y no se envía nada.

## Fuera
- Nada.

## Notas
- Ninguna.
