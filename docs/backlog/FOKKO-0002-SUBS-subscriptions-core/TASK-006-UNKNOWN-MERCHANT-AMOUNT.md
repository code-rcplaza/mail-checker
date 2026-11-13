# FOKKO-0002-SUBS/TASK-006 · Monto de un comercio desconocido
Feature: `internal/subscriptions/domain`

Como usuario quiero que un cobro repetido de un comercio desconocido cuente como
suscripción solo si su monto es estable (o cambió de plan una vez).

## Criterios
- Dados montos 20, 20, 21, cuando se analiza, entonces es consistente, monto actual 21,
  sin cambio de plan.
- Dados montos 20, 20, 100, 100, cuando se analiza, entonces es consistente, monto
  actual 100 y un cambio de plan 20 → 100.
- Dados montos 8.000, 14.500, 5.200, cuando se analiza, entonces no es consistente
  (no es suscripción).

## Notas
- Regla: se comparan cobros consecutivos; con a lo sumo un salto > 25 % es
  consistente (ese salto es el cambio de plan); con más de uno, no.
