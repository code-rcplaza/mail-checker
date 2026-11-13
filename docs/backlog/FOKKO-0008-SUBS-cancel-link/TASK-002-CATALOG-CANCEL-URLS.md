# FOKKO-0008-SUBS/TASK-002 · Enlaces de los servicios del catálogo
Feature: `internal/subscriptions/adapters/catalog`

Como usuario quiero que Netflix, Spotify y el resto muestren su enlace de cancelación, para cancelar en un paso.

## Criterios
- Dado el descriptor de cada servicio con enlace verificado, cuando se llama `Identify`, entonces devuelve su enlace.
- Dado un servicio sin enlace verificado, cuando se llama `Identify`, entonces el enlace es vacío.
- Dado cualquier servicio del catálogo, entonces su enlace, si existe, es `https`.

## Fuera
- Nada.

## Notas
- Hacen falta los enlaces reales verificados por el usuario; no se inventan.
