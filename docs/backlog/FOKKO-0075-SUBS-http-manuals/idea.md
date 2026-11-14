# FOKKO-0075-SUBS · Manuales

## Problema
Agregar permite crear, editar y quitar suscripciones manuales.

## Depende de
- FOKKO-0032-SUBS y FOKKO-0064-SUBS.

## Alcance
- `POST /api/subscriptions/manual` y `PUT`/`DELETE /api/subscriptions/manual/{publicId}`.

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 `POST /api/subscriptions/manual`](TASK-001-MANUAL-CREATE-HANDLER.md)
- [x] [TASK-002 `PUT`/`DELETE` de una manual](TASK-002-MANUAL-EDIT-DELETE-HANDLERS.md)
