# FOKKO-0072D-SUBS/TASK-001 · Total del mes
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver mi gasto mensual en pesos, para tener una referencia.

## Criterios
- Dado un total con dólares convertidos, cuando se serializa, entonces trae el monto en CLP y la marca de aproximado.
- Dado un total sin tipo de cambio, cuando se serializa, entonces trae la marca de incompleto.
- Dado un total solo en pesos, cuando se serializa, entonces no trae ninguna de las dos marcas.

## Fuera
- Nada.

## Notas
- Ninguna.
