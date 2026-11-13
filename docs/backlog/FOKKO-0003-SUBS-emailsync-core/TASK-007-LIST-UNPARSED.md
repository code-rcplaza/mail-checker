# FOKKO-0003-SUBS/TASK-007 · `ListUnparsed`
Feature: `internal/emailsync/usecases`

Como usuario quiero ver los correos que la app no reconoció para saber qué avisos se
están perdiendo.

## Criterios
- Dados correos no reconocidos guardados, cuando se llama `ListUnparsed(ctx, 50)`,
  entonces devuelve los del repositorio con ese límite.
- Dado `limit` = 0 o 101, cuando se llama `ListUnparsed`, entonces `ErrInvalid` con
  `emailsync.limit_out_of_range` y no se consulta el repositorio.
- Dado un error del repositorio, cuando se llama `ListUnparsed`, entonces devuelve el
  error envuelto con contexto.

## Fuera
- Endpoint `/api/unparsed` → fase 5.

## Notas
- Rango 1–100 como constantes con nombre.
