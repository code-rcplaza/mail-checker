# FOKKO-0072B-SUBS/TASK-001 · Suscripciones en el JSON
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver mis suscripciones con su monto y próximo cobro, para revisarlas.

## Criterios
- Dada una suscripción, cuando se serializa, entonces trae `lineKey`, nombre o alias, ciclo, próximo cobro, banco y tarjeta, y las marcas de alza y de "por cancelar".
- Dada una suscripción en USD, cuando se serializa, entonces trae el monto original y el equivalente en CLP con la marca de aproximado.

## Fuera
- Nada.

## Notas
- Ninguna.
