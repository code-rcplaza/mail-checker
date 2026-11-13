# FOKKO-0102-SUBS/TASK-007 · Enviar una decisión por HTTP
Feature: `frontend/src/features/subscriptions/infrastructure`



## Criterios
- Dada una decisión, cuando se llama el gateway, entonces envía `POST /api/subscriptions/decisions` con la clave de la línea.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
