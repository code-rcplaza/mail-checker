# 007-06 · Resumen por línea con los escenarios reales
Feature: `internal/subscriptions/domain`

Como usuario quiero que el resumen muestre cada suscripción paralela y sume bien mi
gasto mensual.

## Criterios
- Dos Netflix (mismo precio $12.990 o distinto $12.990 / $7.990; misma o distinta
  cuenta) → dos suscripciones activas y gasto = suma de ambas.
- Dos Apple de $5.990 (días 5 y 9) → dos suscripciones activas, gasto $11.980.
- Dos pólizas `SEGUROS CLC` el mismo día con montos distintos → dos posibles (no
  "otros cobros").
- Google con $21.700 y $11.000 mensuales → dos activas, gasto $32.700.
- Netflix que sube de $12.990 a $17.990 en la misma cuenta → una activa con monto
  $17.990 y un cambio de plan.
- Un crédito suelto de Claude USD 5 junto a la línea de USD 20 → la línea activa con
  USD 20 y el USD 5 en otros cobros.
- Cada `Subscription` lleva su padre y su clave de línea; una decisión guardada con la
  clave de una línea posible se aplica solo a esa línea.

## Fuera
- Alias → 007-07.

## Notas
- Escenarios sacados de la auditoría con cartolas reales; cobros sintéticos.
