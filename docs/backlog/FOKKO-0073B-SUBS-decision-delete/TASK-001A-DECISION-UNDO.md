# FOKKO-0073B-SUBS/TASK-001A · Deshacer una decisión
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dado `?groupKey=` con una decisión guardada, cuando se envía, entonces responde 204 y la línea vuelve a posibles.
- Dada una clave sin decisión, cuando se envía, entonces responde 204 (idempotente).

## Fuera
- Nada.

## Notas
- Ninguna.
