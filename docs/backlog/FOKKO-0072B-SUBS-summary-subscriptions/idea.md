# FOKKO-0072B-SUBS · `GET /api/summary`: suscripciones

## Problema
El resumen debe listar las suscripciones del usuario con lo que las pantallas muestran de cada una.

## Depende de
- FOKKO-0072A-SUBS, FOKKO-0023-SUBS y FOKKO-0030-SUBS.

## Alcance
- Forma JSON de las suscripciones activas y confirmadas.

## Fuera
- Posibles → FOKKO-0072C-SUBS.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Suscripciones en el JSON](TASK-001-SUMMARY-DTO-SUBSCRIPTIONS.md)
- [x] [TASK-002 Sin datos internos](TASK-002-SUMMARY-DTO-NO-INTERNALS.md)
