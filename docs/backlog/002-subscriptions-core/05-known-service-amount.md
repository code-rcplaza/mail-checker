# 002-05 · Monto de un servicio conocido
Feature: `internal/subscriptions/domain`

Como usuario quiero ver el monto actual de una suscripción conocida y sus cambios de
plan, sin que cobros sueltos del mismo servicio la ensucien.

## Criterios
- Dados cobros mensuales de USD 20, 20, 100, 100, cuando se analiza el monto, entonces
  el monto actual es USD 100 y hay un cambio de plan 20 → 100 en la fecha del primer 100.
- Dados cobros de 6.590, 6.590, 7.190, cuando se analiza, entonces el monto actual es
  7.190 y no hay cambio de plan (variación ≤ 25 %).
- Dados cobros mensuales de USD 20 y un cobro extra de USD 5 a 10 días de uno de
  ellos, cuando se analiza, entonces el extra va a otros cobros del servicio y no
  cambia el monto actual.

## Fuera
- Comercio desconocido → 002-06.
