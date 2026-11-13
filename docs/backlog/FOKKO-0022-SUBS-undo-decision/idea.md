# FOKKO-0022-SUBS · Deshacer una decisión

## Problema
El usuario puede confirmar o descartar una posible por error. Está en el diseño pero no en el código (`SubscriptionDecisions` no tiene `Delete`).

## Depende de
- FOKKO-0021-SUBS.

## Alcance
- `SubscriptionDecisions.Delete` y el caso de uso `UndoDecision`.

## Fuera
- El endpoint `DELETE /api/subscriptions/decisions` → FOKKO-0073A-SUBS y FOKKO-0073B-SUBS.

## Features afectadas
- `internal/subscriptions` (`ports`, `usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 `Delete` en el puerto de decisiones](TASK-001-DECISIONS-DELETE.md)
- [x] [TASK-002A Deshacer una decisión](TASK-002A-UNDO-DECISION.md)
- [x] [TASK-002B No tocar decisiones ajenas](TASK-002B-UNDO-OTHER-USER.md)
