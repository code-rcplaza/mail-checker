# FOKKO-0040-SUBS/TASK-002B · No duplicar un correo registrado
Feature: `internal/auth/usecases`



## Criterios
- Dado un correo ya registrado, cuando se crea, entonces devuelve el usuario existente sin error y no guarda otro (respuesta neutra).
- Dado el mismo correo en mayúsculas distintas, cuando se crea, entonces cuenta como el mismo.

## Fuera
- Nada.

## Notas
- Ninguna.
