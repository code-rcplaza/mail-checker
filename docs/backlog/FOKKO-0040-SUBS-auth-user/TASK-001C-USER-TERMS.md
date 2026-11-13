# FOKKO-0040-SUBS/TASK-001C · Exigir los términos
Feature: `internal/auth/domain`



## Criterios
- Dados los términos sin aceptar, cuando se construye, entonces devuelve `ErrInvalid` con `auth.terms_required`.

## Fuera
- Nada.

## Notas
- Las llaves `auth.*` son constantes en `auth/domain`.
