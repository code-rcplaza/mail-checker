# FOKKO-0076-SUBS/TASK-001 · `POST /api/subscriptions/price-alerts/decision`
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero decir que ya vi el alza, para que deje de avisarme.

## Criterios
- Dada una línea con alza, cuando se envía `{lineKey}`, entonces responde 204.
- Dada una línea sin alza vigente, cuando se envía, entonces responde 404.
- Dado un cuerpo con el campo `later`, cuando se envía, entonces responde 400 con `common.invalid_json` (campo desconocido).

## Fuera
- Nada.

## Notas
- Ninguna.
