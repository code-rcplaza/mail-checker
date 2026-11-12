# 002-13 · Armar el resumen
Feature: `internal/subscriptions/domain`

Como usuario quiero que cada grupo de cobros quede en la sección que le corresponde
(suscripciones, posibles, inactivas u otros cobros) con mi gasto mensual, respetando
lo que decidí sobre las posibles.

## Criterios
- Dado un servicio conocido activo, cuando se arma el resumen, entonces está en
  suscripciones con su monto actual; si está inactivo, en inactivas.
- Dada una posible sin decisión, cuando se arma, entonces está en posibles y no suma
  al gasto.
- Dada una posible confirmada, cuando se arma, entonces está en suscripciones con la
  cadencia que eligió el usuario y suma al gasto.
- Dada una posible descartada, cuando se arma, entonces sus cobros están en otros
  cobros.
- Dado un grupo que no es suscripción, cuando se arma, entonces sus cobros están en
  otros cobros.
- Dadas suscripciones activas, cuando se arma, entonces el gasto mensual es el de
  `MonthlySpend` sobre ellas.

## Fuera
- Puertos, `months`, lectura y deduplicación → 002-09.
- Cobros sueltos de un servicio conocido (ej. créditos de API): sin sección propia
  hasta que se arme la vista.

## Notas
- `BuildSummary(groups, decisions, looksLikeSubscription, now) Summary`, función pura.
- `Summary{Subscriptions, Possible, Inactive, OtherCharges, MonthlySpend}`.
- `Subscription` gana `Charges` para llevar a otros cobros los de los grupos que no
  son suscripción o fueron descartados.
- Monto actual: `AnalyzeKnownServiceAmount` (conocido) o
  `AnalyzeUnknownMerchantAmount` (desconocido). `Detection` no cambia.
