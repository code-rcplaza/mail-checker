# FOKKO-0008-SUBS/TASK-001 · `Service` con enlace de cancelación
Feature: `internal/subscriptions/domain`

Como usuario quiero que cada servicio conocido traiga el enlace para cancelarlo, para llegar al sitio correcto sin buscarlo.

## Criterios
- Dado un `Service` con enlace, cuando se lee, entonces devuelve ese enlace.
- Dado un `Service` sin enlace, cuando se lee, entonces el enlace es vacío y no es un error.
- Dado un enlace que no es `https`, cuando se construye el servicio, entonces se rechaza con `ErrInvalid`.

## Fuera
- Los enlaces concretos de cada servicio → TASK-002.

## Notas
- Un enlace `http` o con esquema raro nunca debe salir por la API.
