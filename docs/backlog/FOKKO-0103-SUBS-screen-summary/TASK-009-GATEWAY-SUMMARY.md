# FOKKO-0103-SUBS/TASK-009 · Leer el resumen por HTTP
Feature: `frontend/src/features/subscriptions/infrastructure`



## Criterios
- Dado un 200, cuando se llama el gateway, entonces mapea el DTO al dominio con montos como `Money`.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
