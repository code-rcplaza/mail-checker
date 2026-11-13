# FOKKO-0095-SUBS/TASK-003 · Iniciar sesión por HTTP
Feature: `frontend/src/features/auth/infrastructure`



## Criterios
- Dado `POST /api/auth/login` con 204, cuando se llama el gateway, entonces resuelve vacío.
- Dado un 400, entonces lanza `ApiError` con la llave.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
