# FOKKO-0077B-SUBS/TASK-001 · `POST /api/subscriptions/cancellations`
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero marcar una suscripción como "por cancelar", para recordarlo.

## Criterios
- Dado `{lineKey}` existente, cuando se envía, entonces responde 204.
- Dada una línea inexistente, cuando se envía, entonces responde 404 con `subscriptions.line_not_found`.

## Fuera
- Nada.

## Notas
- Ninguna.
