# FOKKO-0111-SUBS/TASK-005 · Eliminar todo por HTTP
Feature: `frontend/src/features/profile/infrastructure`



## Criterios
- Dado `DELETE /api/me/data` con 204, cuando se llama el gateway, entonces resuelve vacío.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
