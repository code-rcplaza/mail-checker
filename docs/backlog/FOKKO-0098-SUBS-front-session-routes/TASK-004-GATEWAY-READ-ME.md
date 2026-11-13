# FOKKO-0098-SUBS/TASK-004 · Leer la sesión por HTTP
Feature: `frontend/src/features/auth/infrastructure`



## Criterios
- Dado un 200, cuando se llama el gateway, entonces devuelve el usuario mapeado.
- Dado un 401, entonces devuelve `null`.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
