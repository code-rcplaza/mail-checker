# FOKKO-0020A-SUBS · `UserID` en `Charge`

## Problema
El producto es multiusuario, pero `Charge` y los puertos de `statements` no conocen al dueño. Sin esto un usuario podría ver las cartolas de otro.

Esta idea cubre solo: `UserID` en `Charge`. El resto del trabajo original está en las ideas hermanas `FOKKO-0020*-SUBS`.

## Alcance
- `UserID` en `Charge`

## Fuera
- `subscriptions` → FOKKO-0021-SUBS.
- Cómo se obtiene el usuario en HTTP → FOKKO-0064-SUBS.

## Features afectadas
- `internal/shared` (`domain`)
- `internal/statements` (`ports`, `usecases`)

## Dudas abiertas
- Cómo viaja el `userID` por las capas: se propone un parámetro explícito después de `ctx` (claro y testeable) en lugar de ir dentro del contexto.
- `emailsync` y `accounts` (post-MVP) también usan `Charge`: solo se adapta lo que deja de compilar; sus puertos no ganan `userID` hasta el post-MVP.

## Tareas
- [ ] [TASK-001 `UserID` en `Charge`](TASK-001-CHARGE-USER-ID.md)
