# FOKKO-0040-SUBS/TASK-001B · Rechazar nombre o correo inválido
Feature: `internal/auth/domain`



## Criterios
- Dado un nombre vacío o de solo espacios, cuando se construye, entonces devuelve `ErrInvalid` con `auth.name_required`.
- Dado un nombre de más de 60 caracteres, cuando se construye, entonces devuelve `ErrInvalid` con `auth.name_too_long`.
- Dado un correo con formato inválido, cuando se construye, entonces devuelve `ErrInvalid` con `auth.email_invalid`.

## Fuera
- Nada.

## Notas
- Las llaves `auth.*` son constantes en `auth/domain`.
