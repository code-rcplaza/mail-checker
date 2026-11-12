# 002-09 · `GetSummary`
Feature: `internal/subscriptions/usecases` + `ports`

Como usuario quiero un resumen del período elegido en cuatro secciones:
suscripciones, posibles, inactivas y otros cobros, con mi gasto mensual.

## Criterios
- Dados un cobro por correo y el mismo en la cartola, cuando se pide el resumen,
  entonces la suscripción cuenta ese cobro una sola vez.
- Dado `months = 6`, cuando se pide el resumen, entonces el `ChargeReader` se consulta
  desde hoy menos 6 meses.
- Dado un servicio conocido activo, cuando se pide, entonces está en suscripciones; si
  está inactivo, en inactivas.
- Dada una posible sin decisión, cuando se pide, entonces está en posibles y no suma al
  gasto.
- Dada una posible confirmada, cuando se pide, entonces está en suscripciones con la
  cadencia que eligió el usuario y suma al gasto; si fue descartada, está en otros cobros.
- Dados cobros que no forman suscripción, cuando se pide, entonces están en otros cobros.
- Dado `months` fuera de 1–24, cuando se pide, entonces `ErrInvalid` con
  `subscriptions.months_out_of_range`.
- Dado que el `ChargeReader` o `SubscriptionDecisions` falla, cuando se pide, entonces
  devuelve ese error envuelto.

## Notas
- Puertos `ChargeReader`, `ServiceCatalog` y `SubscriptionDecisions` con dobles en
  memoria; reloj inyectado.
