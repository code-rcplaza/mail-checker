# FOKKO-0015-SUBS · Derivar el tipo de cambio de los traspasos

## Problema
Con el par de traspasos se calcula el tipo de cambio vendedor del banco, que es lo que el usuario paga de verdad.

## Depende de
- FOKKO-0014-SUBS.

## Alcance
- Función pura que calcula el tipo de cambio de un par.

## Fuera
- Elegir cuál usar → FOKKO-0016-SUBS.

## Features afectadas
- `internal/shared` o `internal/statements` (`domain`), según la duda

## Dudas abiertas
- Dónde vive el tipo `ExchangeRate`: lo usan `statements` (lo deriva) y `subscriptions` (lo lee). Se propone `shared/domain`, igual que `Money`.

## Tareas
- [ ] [TASK-001A Calcular el tipo de cambio](TASK-001A-DERIVE-RATE.md)
- [ ] [TASK-001B Rechazar montos inválidos](TASK-001B-DERIVE-RATE-GUARDS.md)
