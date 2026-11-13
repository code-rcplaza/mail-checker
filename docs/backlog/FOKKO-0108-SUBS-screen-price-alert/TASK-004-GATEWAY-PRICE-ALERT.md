# FOKKO-0108-SUBS/TASK-004 · Enviar la decisión de alza por HTTP
Feature: `frontend/src/features/subscriptions/infrastructure`



## Criterios
- Dado `keep`, cuando se llama el gateway, entonces envía `POST /api/price-alerts/decision` con `{lineKey}`.

## Fuera
- Los componentes ya existen en otras tareas de esta idea.

## Notas
- Implementa el puerto de `application` con `lib/api`; el mapper DTO→dominio vive aquí.
