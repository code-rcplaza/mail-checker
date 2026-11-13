# FOKKO-0026-SUBS/TASK-001 · Detectar el alza entre dos cobros
Feature: `internal/subscriptions/domain`

Como usuario quiero que la app me diga cuando un cobro subió, para enterarme sin revisar la cartola.

## Criterios
- Dados dos cobros consecutivos de CLP 10.000 y CLP 10.050, cuando se detecta, entonces hay alza con antes 10.000, después 10.050 y la fecha del cobro nuevo.
- Dados dos cobros iguales, cuando se detecta, entonces no hay alza.
- Dado un cobro menor que el anterior, cuando se detecta, entonces no hay alza.

## Fuera
- Nada.

## Notas
- Ninguna.
