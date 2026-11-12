# 002-08 · Gasto mensual
Feature: `internal/subscriptions/domain`

Como usuario quiero saber cuánto gasto al mes en suscripciones, por moneda.

## Criterios
- Dadas dos suscripciones activas mensuales de 6.590 y 9.990 CLP y una de USD 20,
  cuando se suma, entonces CLP 16.580 y USD 20 por separado.
- Dada una activa anual de 120.000 CLP, cuando se suma, entonces aporta 10.000 al mes.
- Dada una anual cuyo monto / 12 no es exacto, cuando se suma, entonces se redondea a
  la centésima más cercana (mitad hacia arriba).
- Dada una posible confirmada como mensual, cuando se suma, entonces aporta su monto;
  confirmada como anual (ej. USD 120), aporta monto / 12 (USD 10).
- Dadas suscripciones inactivas, posibles sin decidir o posibles descartadas, cuando se
  suma, entonces no aportan.
