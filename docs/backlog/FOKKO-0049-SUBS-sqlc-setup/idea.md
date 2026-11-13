# FOKKO-0049-SUBS · `sqlc.yaml` y generación de queries (la migración base no crea tablas)

## Problema
Las consultas SQL viven solo en queries de sqlc (principio DRY del proyecto), pero no existe `sqlc.yaml` y la migración base no crea tablas. Sin esto ningún repo se puede escribir.

## Alcance
- `backend/sqlc.yaml`: motor SQLite, esquema tomado de las migraciones de goose, una entrada por feature con sus `queries.sql` y su salida generada.
- El comando `sqlc` entra en "Comandos" de `CLAUDE.md`.

## Fuera
- Las tablas y consultas de cada feature → una idea de repo cada una (FOKKO-0050-SUBS en adelante).

## Features afectadas
- `backend` (herramientas)

## Dudas abiertas
- Se propone commitear el código generado, para compilar sin instalar `sqlc`. ¿Se acepta?
- Editar `CLAUDE.md` (sección Comandos) pide tu ok al hacerlo.
- Versión de `sqlc` a fijar: se decide al escribir la tarea.

## Tareas
- [ ] [TASK-001 Configuración de sqlc](TASK-001-SQLC-CONFIG.md)
- [ ] [TASK-002 Comprobar que lo generado está al día](TASK-002-SQLC-VERIFY.md)
