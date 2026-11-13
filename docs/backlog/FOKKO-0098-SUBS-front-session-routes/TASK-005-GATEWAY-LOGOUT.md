# FOKKO-0098-SUBS/TASK-005 · Cerrar sesión por HTTP
Feature: `frontend/src/features/auth/infrastructure`



## Criterios
- Dado `POST /api/auth/logout` con 204, cuando se llama el gateway, entonces resuelve vacío.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
