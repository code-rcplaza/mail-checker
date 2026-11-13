# FOKKO-0034-SUBS/TASK-001 · Ocultar una línea
Feature: `internal/subscriptions/usecases`

Como usuario quiero dejar de seguir una suscripción que ya no me importa, para que no moleste.

## Criterios
- Dada una línea existente, cuando se deja de seguir, entonces se guarda la marca por clave de línea.
- Dada una línea ya oculta, cuando se deja de seguir, entonces no falla (idempotente).
- Dada una línea inexistente o de otro usuario, cuando se deja de seguir, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
