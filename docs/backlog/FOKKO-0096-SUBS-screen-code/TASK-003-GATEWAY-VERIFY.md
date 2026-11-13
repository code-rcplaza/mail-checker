# FOKKO-0096-SUBS/TASK-003 · Verificar el código por HTTP
Feature: `frontend/src/features/auth/infrastructure`



## Criterios
- Dado un código correcto, cuando se llama el gateway, entonces resuelve vacío.
- Dado un 400, entonces lanza `ApiError`.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
