# FOKKO-0094-SUBS/TASK-003 · Registrar por HTTP
Feature: `frontend/src/features/auth/infrastructure`



## Criterios
- Dado `POST /api/auth/register` con 204, cuando se llama el gateway, entonces resuelve vacío.
- Dado un 400 o 409, entonces lanza `ApiError` con la llave y los campos.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
