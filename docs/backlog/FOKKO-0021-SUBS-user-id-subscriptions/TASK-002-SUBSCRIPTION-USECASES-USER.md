# FOKKO-0021-SUBS/TASK-002 · Casos de uso con usuario
Feature: `internal/subscriptions/usecases`

Como usuario quiero que mi resumen use solo mis cobros, para no ver los de otra persona.

## Criterios
- Dado `GetSummary` con un usuario, cuando se ejecuta, entonces solo lee los cargos, decisiones y alias de ese usuario.
- Dado `DecidePossibleSubscription` con la línea de otro usuario, cuando se ejecuta, entonces devuelve `ErrNotFound`.
- Dado `NameLine` con la línea de otro usuario, cuando se ejecuta, entonces devuelve `ErrNotFound` con `subscriptions.line_not_found`.

## Fuera
- Nada.

## Notas
- Ninguna.
