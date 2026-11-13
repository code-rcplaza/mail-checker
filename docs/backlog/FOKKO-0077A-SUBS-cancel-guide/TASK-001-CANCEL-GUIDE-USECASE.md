# FOKKO-0077A-SUBS/TASK-001 · Buscar el enlace de una línea
Feature: `internal/subscriptions/usecases`

Como usuario quiero ir directo al sitio para cancelar, para no buscarlo.

## Criterios
- Dada una línea de un servicio con enlace, cuando se pide la guía, entonces devuelve su enlace.
- Dada una línea de un servicio o comercio sin enlace, cuando se pide la guía, entonces devuelve vacío sin error.
- Dada una línea inexistente o de otro usuario, cuando se pide la guía, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
