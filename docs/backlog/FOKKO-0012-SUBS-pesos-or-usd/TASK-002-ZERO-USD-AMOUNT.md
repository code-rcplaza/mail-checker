# FOKKO-0012-SUBS/TASK-002 · US$ en cero
Feature: `internal/statements/domain`

Como usuario quiero que un cobro con US$ 0,00 no rompa la lectura, para no perder el resto de la cartola.

## Criterios
- Dado origen mayor que cero y US$ 0,00, cuando se decide, entonces es en pesos (la razón es infinita).
- Dados origen y US$ ambos en cero, cuando se decide, entonces no es en pesos y no hay división por cero.

## Fuera
- Nada.

## Notas
- Ninguna.
