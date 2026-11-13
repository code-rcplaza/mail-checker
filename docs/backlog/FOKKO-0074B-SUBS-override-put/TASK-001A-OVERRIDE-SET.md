# FOKKO-0074B-SUBS/TASK-001A · Guardar o borrar el override
Feature: `internal/subscriptions/adapters/http`



## Criterios
- Dado `{lineKey, amount, billingDay}` válido, cuando se envía, entonces responde 204.
- Dados monto y día ausentes, cuando se envía, entonces se borra el override y responde 204.

## Fuera
- Nada.

## Notas
- Ninguna.
