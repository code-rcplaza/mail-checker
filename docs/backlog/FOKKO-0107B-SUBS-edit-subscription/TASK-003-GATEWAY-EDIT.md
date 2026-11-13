# FOKKO-0107B-SUBS/TASK-003 · Editar por HTTP
Feature: `frontend/src/features/subscriptions/infrastructure`



## Criterios
- Dado un cambio válido, cuando se llama el gateway, entonces envía la edición.
- Dado un 400, entonces lanza `ApiError` con los campos.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
