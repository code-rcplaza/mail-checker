# 002-09 · `GetSummary`
Feature: `internal/subscriptions/usecases` + `ports`

Como usuario quiero un resumen del período elegido en cuatro secciones:
suscripciones, posibles, inactivas y otros cobros, con mi gasto mensual.

## Criterios
- Dados un cobro por correo y el mismo en la cartola, cuando se pide el resumen,
  entonces la suscripción cuenta ese cobro una sola vez.
- Dado `months = 6`, cuando se pide el resumen, entonces el `ChargeReader` se consulta
  desde hoy menos 6 meses.
- Dados cobros de un servicio conocido activo y una posible confirmada guardada en
  `SubscriptionDecisions`, cuando se pide, entonces ambos están en suscripciones (las
  secciones salen de `BuildSummary`, ver 002-13).
- Dado `months` fuera de 1–24, cuando se pide, entonces `ErrInvalid` con
  `subscriptions.months_out_of_range`.
- Dado que el `ChargeReader` o `SubscriptionDecisions` falla, cuando se pide, entonces
  devuelve ese error envuelto.

## Fuera
- Reparto en secciones y gasto mensual → 002-13.

## Notas
- Requiere 002-13.
- Puertos `ChargeReader`, `ServiceCatalog` y `SubscriptionDecisions` (solo `List`; `Save`
  llega con 002-11) con dobles en memoria; reloj inyectado.
