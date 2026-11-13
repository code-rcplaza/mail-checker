# FOKKO-0008-SUBS · Enlace de cancelación en el catálogo

## Problema
La pantalla "Cómo cancelar" necesita llevar al usuario al sitio del servicio. Ese dato es del catálogo (mismo para todos los bancos) y hoy `Service` solo trae nombre y si es umbrella.

## Depende de
- FOKKO-0006-SUBS (la tabla del catálogo).

## Alcance
- `Service` gana el enlace para cancelar o administrar la suscripción.
- Cada servicio del catálogo con enlace verificado lo devuelve desde `Identify`.

## Fuera
- El endpoint `cancel-guide` → FOKKO-0077A-SUBS y FOKKO-0077B-SUBS.
- Los 3 pasos de cancelación (texto fijo del front) → FOKKO-0109-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`, `adapters/catalog`)

## Dudas abiertas
- Los enlaces de los 10 servicios del catálogo no se inventan: se piden y se verifican uno a uno.
- Para un servicio umbrella (App Store, Google, Microsoft) el enlace es el de administrar suscripciones de la tienda, no el de un servicio concreto. ¿Se acepta?

## Tareas
- [ ] [TASK-001 `Service` con enlace de cancelación](TASK-001-CANCEL-URL-FIELD.md)
- [ ] [TASK-002 Enlaces de los servicios del catálogo](TASK-002-CATALOG-CANCEL-URLS.md)
