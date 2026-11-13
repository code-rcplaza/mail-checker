# FOKKO-0013-SUBS/TASK-002 · Guardar en pesos o en dólares
Feature: `internal/statements/adapters/bci`

Como usuario quiero que un cobro internacional se guarde en la moneda correcta, para que su línea coincida con la del estado nacional.

## Criterios
- Dada una línea cuyo origen es en pesos según FOKKO-0012-SUBS, cuando se parsea, entonces el cargo queda en CLP con el monto de origen.
- Dada una línea que no es en pesos, cuando se parsea, entonces el cargo queda en USD con el monto en US$.
- Dada una línea en euros con su monto en US$, cuando se parsea, entonces el cargo queda en USD (test sintético).

## Fuera
- Nada.

## Notas
- El caso en euros usa un ejemplo sintético hasta tener uno real.
