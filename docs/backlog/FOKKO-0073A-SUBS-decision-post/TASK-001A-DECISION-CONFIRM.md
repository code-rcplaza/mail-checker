# FOKKO-0073A-SUBS/TASK-001A · Confirmar con cadencia
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dado `{groupKey, decision:"confirm", cadence:"monthly"}`, cuando se envía, entonces responde 204 y la línea queda confirmada.
- Dado `confirm` sin cadencia, cuando se envía, entonces responde 400 con `subscriptions.cadence_required`.

## Fuera
- Nada.

## Notas
- Ninguna.
