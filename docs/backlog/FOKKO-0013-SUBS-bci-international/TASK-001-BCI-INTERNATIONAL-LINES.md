# FOKKO-0013-SUBS/TASK-001 · Leer las líneas del estado internacional
Feature: `internal/statements/adapters/bci`

Como usuario quiero que mis compras en dólares se lean de la cartola, para ver mis suscripciones extranjeras.

## Criterios
- Dada una línea internacional, cuando se parsea, entonces trae descriptor, fecha, monto de origen y monto en US$.
- Dado `US$ 20,00`, cuando se parsea, entonces vale `2000` centésimas.
- Dadas varias líneas, cuando se parsea, entonces cada una trae su número de línea.

## Fuera
- Nada.

## Notas
- Ninguna.
