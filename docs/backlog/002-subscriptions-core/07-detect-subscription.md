# 002-07 · Detectar una suscripción
Feature: `internal/subscriptions/domain`

Como usuario quiero saber si cada grupo de cobros es una suscripción (activa o
inactiva), una posible suscripción, o no es suscripción.

## Criterios
- Dado un servicio conocido mensual con último cobro hace 10 días, cuando se detecta,
  entonces es suscripción activa mensual.
- Dado un servicio mensual con último cobro hace 50 días, cuando se detecta, entonces
  es suscripción inactiva.
- Dado un servicio anual con último cobro hace 300 días, cuando se detecta, entonces
  es suscripción activa anual.
- Dado un servicio conocido con un solo cobro, cuando se detecta, entonces es posible
  (no suscripción): con 1 cobro no se sabe si es mensual o anual.
- Dado un comercio desconocido con ≥ 2 cobros mensuales de monto consistente, cuando
  se detecta, entonces es posible.
- Dado un comercio desconocido con 1 cobro cuyo descriptor parece suscripción
  ("… SUBSCRIPTION", "MEMBERSHIP"), cuando se detecta, entonces es posible.
- Dado un comercio desconocido con 1 cobro sin descriptor de suscripción, o con montos
  no consistentes, cuando se detecta, entonces no es suscripción.
- Dada una posible cuyo último cobro está fuera de su ventana de actividad (45 días
  mensual, 400 anual, 365 sin cadencia), cuando se detecta, entonces es inactiva; dentro
  de la ventana sigue siendo posible (ej. conocido con 1 cobro hace 300 días).

## Notas
- "Hoy" llega como parámetro (`now time.Time`), nunca `time.Now()` en `domain`.
- `looksLikeSubscription` llega como función (igual que `identify`), desde el catálogo.
