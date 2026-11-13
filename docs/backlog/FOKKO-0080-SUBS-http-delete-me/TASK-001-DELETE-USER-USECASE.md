# FOKKO-0080-SUBS/TASK-001 · Eliminar al usuario
Feature: `internal/auth/usecases`

Como usuario quiero borrar mi cuenta y mis datos, para dejar de existir en la app.

## Criterios
- Dado un usuario, cuando se elimina, entonces el repositorio lo borra y se cierran sus sesiones.
- Dado un usuario inexistente, cuando se elimina, entonces devuelve `ErrNotFound`.
- Dado un fallo del repositorio, cuando se elimina, entonces devuelve el error con contexto.

## Fuera
- Nada.

## Notas
- Ninguna.
