# FOKKO-0053B-SUBS · Repo de cartolas: lectura y borrado

## Problema
El usuario debe poder ver y borrar solo sus cartolas, y borrar una cartola borra sus cargos.

## Depende de
- FOKKO-0053A-SUBS.

## Alcance
- `List` y `Delete` del repositorio de cartolas.

## Fuera
- Los endpoints de listar y borrar cartolas son post-MVP.

## Features afectadas
- `internal/statements` (`adapters/sqlite`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Listar mis cartolas](TASK-001A-STATEMENTS-LIST.md)
- [x] [TASK-001B Borrar una cartola con sus cargos](TASK-001B-STATEMENTS-DELETE.md)
