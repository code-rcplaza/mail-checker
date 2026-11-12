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
    - mínimo 2 cobros que cumplan el intervalo; servicio conocido con 1 cobro → probable;
    - activa: último cobro ≤ 45 días (≤ 400 si es anual) desde hoy (reloj inyectado).
  - **Monto**:
    - servicio conocido: ±25 % no separa; monto actual = último cobro; salto > 25 % =
      cambio de plan; cobros fuera del patrón → otros cobros del servicio;
    - desconocido: compara con el cobro anterior; salto > 25 % que se repite = cambio
      de plan; montos que cambian casi siempre → no es suscripción.
  - **Gasto mensual** de las activas por moneda (anual = monto / 12).
- `subscriptions/ports`: `ChargeReader`, `ServiceCatalog` (dobles en memoria en los tests).
- `subscriptions/usecases`:
  - `GetSummary(months)`: vista unificada (cartola deduplicada contra correo).
  - `ReconcileCharges(months)`: avisados / sin aviso / solo correo.

## Fuera
- Catálogo real de servicios (descriptores, prefijos de procesadores) → fase 3.
- Repositorio SQLite de cobros (`ChargeReader` real) → fase 4.
- Endpoints `/api/summary` y `/api/reconciliation` → fase 5.
- `ParsedCharge` → idea de `emailsync`.
- Tolerancia configurable por servicio → YAGNI.

## Features afectadas
- `internal/shared/domain`
- `internal/subscriptions`

## Dudas abiertas
- Ninguna. Datos de prueba sintéticos con servicios genéricos ("STREAMING CO",
  "MUSIC APP"); sin datos personales.

## Historias
(se completa con /stories)
