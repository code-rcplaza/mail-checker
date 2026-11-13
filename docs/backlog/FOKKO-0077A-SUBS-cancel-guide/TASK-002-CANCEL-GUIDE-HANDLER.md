# FOKKO-0077A-SUBS/TASK-002 · `GET /api/subscriptions/cancel-guide`
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero abrir el sitio del servicio desde la app, para cancelar en un paso.

## Criterios
- Dado `?lineKey=` con enlace, cuando se pide, entonces responde 200 con `{url}`.
- Dado `?lineKey=` sin enlace, cuando se pide, entonces responde 200 con `url` vacío.
- Dado `lineKey` ausente o inexistente, cuando se pide, entonces responde 400 o 404 según corresponda.

## Fuera
- Nada.

## Notas
- Ninguna.
