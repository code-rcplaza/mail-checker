# FOKKO-0107A-SUBS/TASK-003 · Crear la manual por HTTP
Feature: `frontend/src/features/subscriptions/infrastructure`



## Criterios
- Dado un formulario válido, cuando se llama el gateway, entonces envía `POST /api/subscriptions/manual`.
- Dado un 400, entonces lanza `ApiError` con los campos.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
