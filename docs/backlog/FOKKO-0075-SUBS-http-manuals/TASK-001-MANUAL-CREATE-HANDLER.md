# FOKKO-0075-SUBS/TASK-001 · `POST /api/subscriptions/manual`
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero agregar una suscripción que no aparece en mi cartola, para ver todo mi gasto.

## Criterios
- Dado `{name, amount, billingDay}` válido, cuando se envía, entonces responde 201 con `id` igual al `public_id`.
- Dados campos inválidos, cuando se envía, entonces responde 400 con `common.validation_failed` y `fields`.

## Fuera
- Nada.

## Notas
- Ninguna.
