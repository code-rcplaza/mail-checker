# FOKKO-0040-SUBS/TASK-002B · Rechazar un correo repetido
Feature: `internal/auth/usecases`



## Criterios
- Dado un correo ya registrado, cuando se crea, entonces devuelve `ErrConflict` con `auth.email_taken` y no guarda nada.
- Dado el mismo correo en mayúsculas distintas, cuando se crea, entonces cuenta como el mismo.

## Fuera
- Nada.

## Notas
- Ninguna.
