# FOKKO-0023-SUBS/TASK-001 · Próximo cobro mensual
Feature: `internal/subscriptions/domain`

Como usuario quiero saber cuándo me cobran de nuevo, para tener el saldo.

## Criterios
- Dada una línea mensual con último cobro el 1 de septiembre y día típico 1, cuando se calcula, entonces el próximo es el 1 de octubre.
- Dado un último cobro el 3 y día típico 1, cuando se calcula, entonces el próximo cae en el día 1 del mes siguiente (manda el día típico).

## Fuera
- Nada.

## Notas
- Ninguna.
