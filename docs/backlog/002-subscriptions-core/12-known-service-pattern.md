# 002-12 · Patrón de un servicio conocido por respaldo
Feature: `internal/subscriptions/domain`

Como usuario quiero que los cobros sueltos de un servicio conocido se separen bien
aunque lleguen antes del primer cobro, en una suscripción anual o en otra moneda,
para que el monto actual y los cambios de plan sean correctos.

## Criterios
- Dados cobros de USD 5 (día 0), 20 (día 10) y 20 (día 40), cuando se analiza el
  monto, entonces el de 5 va a otros cobros del servicio, el monto actual es USD 20 y
  no hay cambio de plan.
- Dados cobros anuales de USD 120 (día 0) y 120 (día 365) y uno de USD 5 (día 90),
  cuando se analiza, entonces el de 5 va a otros cobros del servicio y los de 120 son
  el patrón.
- Dados cobros mensuales de USD 20, 20 y un último de USD 100 a 30 días (cambio de
  plan reciente, sin pareja todavía), cuando se analiza, entonces el de 100 es del
  patrón, el monto actual es USD 100 y hay un cambio de plan 20 → 100.
- Dados cobros de USD 100 (día 0), 20 (día 30) y 20 (día 60) (baja de plan), cuando se
  analiza, entonces el de 100 es del patrón y hay un cambio de plan 100 → 20.
- Dados cobros mensuales de CLP 5.990, 5.990 y luego USD 7,99, 7,99, cuando se
  analiza, entonces hay un cambio de plan CLP 5.990 → USD 7,99 y el monto actual es
  USD 7,99.
- Dados los cobros de los criterios de 002-05, cuando se analizan, entonces el
  resultado no cambia (sus tests siguen en verde sin tocarlos).

## Fuera
- Comercio desconocido → 002-06.
- Detección de suscripción y estado activa / posible → 002-07.

## Notas
- **Respaldo**: dos cobros se respaldan si tienen la misma moneda, montos a ±25 % y
  un intervalo válido entre ellos (mensual k × 24–38 o anual 330–400). Un cobro con
  respaldo es del patrón.
- **Cadencia del patrón**: `DetectCadence` sobre los cobros con respaldo.
- **Cobro sin respaldo**: es del patrón si está a un intervalo válido para esa
  cadencia del cobro con respaldo más cercano (antes o después); si no, va a otros
  cobros del servicio.
- **Sin ningún respaldo** (1 cobro, o montos todos distintos): se mantiene la regla
  de 002-05 (un cobro a < 24 días del último del patrón es otro cobro).
- **Moneda**: un cambio de moneda entre cobros consecutivos del patrón es cambio de
  plan; `PlanChange` lleva `FromCurrency` y `ToCurrency`.
- Al cerrar, `docs/ARCHITECTURE.md` reemplaza la regla de 24 días por la de respaldo
  y agrega el cambio de moneda.
