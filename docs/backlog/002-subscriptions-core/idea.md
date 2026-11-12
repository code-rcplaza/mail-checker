# 002 · Núcleo de suscripciones (fase 2, parte 1)

## Problema
El armazón está listo, pero no hay lógica: la app todavía no sabe qué es una
suscripción. Esta idea implementa la detección (la parte de más valor) y las dos
vistas sobre los cobros, con TDD y dobles en memoria, sin base ni red.

## Alcance
- `shared/domain`: `Money` (centésimas, int64), `Currency` (`CLP`, `USD`), `Bank`,
  `Card`, `Source` (`email`, `statement`) y `Charge`.
- `subscriptions/domain`:
  - `Service` (nombre, *umbrella*), `Subscription`, `Summary`, `Reconciliation`.
  - **Regla de coincidencia** cartola ↔ correo: mismo banco, mismo servicio, mismo
    monto y ±35 días (una sola función, usada por las dos vistas).
  - **Agrupación**: servicio conocido → por servicio; *umbrella* → por monto exacto +
    moneda; desconocido → por nombre normalizado + moneda.
  - **Recurrencia**:
    - mensual: k × 24–38 días con k ≤ 3; anual: 330–400 días;
    - servicio conocido: ≥ 2 cobros que cumplan el intervalo → suscripción; 1 solo
      cobro → **posible** (el usuario decide);
    - comercio desconocido: **posible** si tiene ≥ 2 cobros consistentes o 1 cobro con
      descriptor de suscripción; el usuario la confirma o la descarta;
    - activa: último cobro ≤ 45 días (≤ 400 si es anual) desde hoy (reloj inyectado).
  - **Monto**:
    - servicio conocido: ±25 % no separa; monto actual = último cobro; salto > 25 % =
      cambio de plan; cobros fuera del patrón → otros cobros del servicio;
    - desconocido: compara con el cobro anterior; salto > 25 % que se repite = cambio
      de plan; montos que cambian casi siempre → no es suscripción.
  - **Gasto mensual** de las activas + posibles confirmadas, por moneda (anual = monto / 12).
- `subscriptions/ports`: `ChargeReader`, `ServiceCatalog`, `SubscriptionDecisions` (dobles
  en memoria en los tests).
- `subscriptions/usecases`:
  - `GetSummary(months)`: vista unificada en cuatro secciones (suscripciones, posibles,
    inactivas, otros cobros), con la cartola deduplicada contra correo.
  - `ReconcileCharges(months)`: avisados / sin aviso / solo correo.
  - `DecidePossibleSubscription(groupKey, confirm(monthly | yearly) | dismiss)`: al confirmar,
    el usuario elige la cadencia.

## Fuera
- Catálogo real de servicios (descriptores, prefijos de procesadores) → fase 3.
- Repositorios SQLite (`ChargeReader`, `SubscriptionDecisions`) → fase 4.
- Endpoints `/api/summary`, `/api/reconciliation` y `/api/subscriptions/decisions` → fase 5.
- `ParsedCharge` → idea de `emailsync`.
- Tolerancia configurable por servicio → YAGNI.

## Features afectadas
- `internal/shared/domain`
- `internal/subscriptions`

## Dudas abiertas
- Ninguna. Datos de prueba sintéticos con servicios genéricos ("STREAMING CO",
  "MUSIC APP"); sin datos personales.

## Historias
- [x] [01 Tipos de cobro](01-charge-types.md)
- [x] [02 Agrupar cobros](02-group-charges.md)
- [x] [03 Regla de coincidencia cartola ↔ correo](03-match-across-sources.md)
- [x] [04 Cadencia](04-cadence.md)
- [x] [05 Monto de un servicio conocido](05-known-service-amount.md)
- [x] [06 Monto de un comercio desconocido](06-unknown-merchant-amount.md)
- [x] [07 Detectar una suscripción](07-detect-subscription.md)
- [x] [08 Gasto mensual](08-monthly-spend.md)
- [ ] [09 `GetSummary`](09-get-summary.md)
- [ ] [10 `ReconcileCharges`](10-reconcile-charges.md)
- [ ] [11 Confirmar o descartar una posible](11-decide-possible-subscription.md)
- [x] [12 Patrón de un servicio conocido por respaldo](12-known-service-pattern.md)
- [x] [13 Armar el resumen](13-build-summary.md)
