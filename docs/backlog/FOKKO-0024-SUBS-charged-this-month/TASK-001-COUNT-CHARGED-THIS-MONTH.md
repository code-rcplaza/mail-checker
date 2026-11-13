# FOKKO-0024-SUBS/TASK-001 · Contar los cobrados del mes
Feature: `internal/subscriptions/domain`

Como usuario quiero saber cuántos cobros ya pasaron este mes, para ver lo que llevo.

## Criterios
- Dadas tres suscripciones con un cobro este mes y una sin, cuando se cuenta, entonces son 3.
- Dado un cobro de las 23:30 del último día del mes anterior en Santiago, cuando se cuenta, entonces no cuenta para este mes.
- Dada una posible sin decidir, cuando se cuenta, entonces no entra.

## Fuera
- Nada.

## Notas
- Ninguna.
