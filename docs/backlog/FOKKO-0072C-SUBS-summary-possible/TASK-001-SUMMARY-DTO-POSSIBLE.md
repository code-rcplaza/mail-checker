# FOKKO-0072C-SUBS/TASK-001 · Posibles con evidencia
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver por qué la app cree que algo es una suscripción, para decidir con información.

## Criterios
- Dada una posible, cuando se serializa, entonces trae su `lineKey`, nombre, monto y los cobros vistos (fecha y monto).
- Dada una posible sin decidir, cuando se serializa, entonces no trae próximo cobro ni cadencia.
- Dada una posible fuera de su ventana (inactiva), cuando se serializa, entonces no aparece.

## Fuera
- Nada.

## Notas
- Ninguna.
