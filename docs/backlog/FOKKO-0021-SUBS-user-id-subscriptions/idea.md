# FOKKO-0021-SUBS · `UserID` en los puertos de `subscriptions`

## Problema
Lo mismo en `subscriptions`: sus puertos de lectura y escritura no filtran por usuario.

## Depende de
- FOKKO-0020A-SUBS y FOKKO-0020B-SUBS.

## Alcance
- `ChargeReader`, `SubscriptionDecisions` y `LineAliases` filtran por usuario.
- `GetSummary`, `DecidePossibleSubscription` y `NameLine` reciben el usuario.

## Fuera
- Los puertos nuevos del MVP nacen con usuario desde su propia idea.

## Features afectadas
- `internal/subscriptions` (`ports`, `usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Puertos de `subscriptions` por usuario](TASK-001-SUBSCRIPTION-PORTS-USER.md)
- [ ] [TASK-002 Casos de uso con usuario](TASK-002-SUBSCRIPTION-USECASES-USER.md)
