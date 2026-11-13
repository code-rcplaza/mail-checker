# FOKKO-0094-SUBS/TASK-001A · Errores del nombre o el correo
Feature: `frontend/src/features/auth/domain`



## Criterios
- Dado un nombre vacío, cuando se valida, entonces el error es `auth.name_required`.
- Dado un correo sin formato válido, cuando se valida, entonces el error es `auth.email_invalid`.

## Fuera
- Nada.

## Notas
- Función pura en `domain`: sin React ni red.
