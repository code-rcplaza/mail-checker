# FOKKO-0026-SUBS/TASK-002 · El cambio de moneda no es alza
Feature: `internal/subscriptions/domain`

Como usuario quiero que pasar de dólares a pesos no cuente como alza, para no ver alertas falsas.

## Criterios
- Dado un cobro en USD seguido de uno en CLP en la misma línea, cuando se detecta, entonces no hay alza (es cambio de plan).
- Dada una línea con un solo cobro, cuando se detecta, entonces no hay alza.

## Fuera
- Nada.

## Notas
- Ninguna.
