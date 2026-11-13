# FOKKO-0099-SUBS/TASK-005 · Subir la cartola por HTTP
Feature: `frontend/src/features/statements/infrastructure`



## Criterios
- Dado un archivo y un banco, cuando se llama el gateway, entonces envía `multipart` a `/api/statements`.
- Dado un 422, entonces lanza `ApiError` con la llave.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
