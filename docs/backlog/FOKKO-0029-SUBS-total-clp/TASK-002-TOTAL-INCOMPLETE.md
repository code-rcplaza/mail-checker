# FOKKO-0029-SUBS/TASK-002 · Total sin tipo de cambio
Feature: `internal/subscriptions/domain`

Como usuario quiero saber si mi total no incluye los dólares, para no confiar en él.

## Criterios
- Dado un gasto en USD y ningún tipo de cambio, cuando se suma, entonces el total solo incluye los pesos y queda marcado incompleto.
- Dado solo gasto en pesos y ningún tipo de cambio, cuando se suma, entonces el total no es incompleto.

## Fuera
- Nada.

## Notas
- Ninguna.
