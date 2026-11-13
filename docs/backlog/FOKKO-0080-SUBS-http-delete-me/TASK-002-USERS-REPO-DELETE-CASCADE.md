# FOKKO-0080-SUBS/TASK-002 · La cascada borra todo
Feature: `internal/auth/adapters/sqlite`

Como usuario quiero que no quede nada mío, para mi privacidad.

## Criterios
- Dado un usuario con datos en todas las tablas, cuando se elimina, entonces no queda ninguna fila suya en ninguna tabla.
- Dado otro usuario con datos, cuando se elimina el primero, entonces los datos del otro quedan intactos.
- Dado el borrado, entonces ocurre en una sola transacción (o todo o nada).

## Fuera
- Nada.

## Notas
- Ninguna.
