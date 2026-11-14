# FOKKO-0044-SUBS · Límite de 3 envíos por hora

## Problema
Sin tope de envíos, alguien podría llenar un correo ajeno de códigos. El límite es de 3 envíos por hora por correo.

## Depende de
- FOKKO-0042-SUBS.

## Alcance
- Política pura de límite y su aplicación al emitir un código.

## Fuera
- El código HTTP 429 de la respuesta → FOKKO-0063-SUBS.

## Features afectadas
- `internal/shared` (`domain`)
- `internal/auth` (`domain`, `usecases`)

## Dudas abiertas
- Hoy los errores compartidos son inválido, no existe y conflicto. Un límite pide un error propio (`ErrTooMany` → HTTP 429). Se propone agregarlo a `shared/domain`.

## Tareas
- [x] [TASK-001 Error de límite excedido](TASK-001-TOO-MANY-ERROR.md)
- [x] [TASK-002 Política de 3 envíos por hora](TASK-002-SEND-LIMIT-POLICY.md)
- [x] [TASK-003 Aplicar el límite al emitir](TASK-003-ENFORCE-SEND-LIMIT.md)
