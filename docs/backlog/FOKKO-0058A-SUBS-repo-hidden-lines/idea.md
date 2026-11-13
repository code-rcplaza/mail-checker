# FOKKO-0058A-SUBS · Repo de líneas ocultas

## Problema
"Dejar de seguir" oculta una línea y la marca debe persistir por clave de línea.

## Depende de
- FOKKO-0034-SUBS y FOKKO-0049-SUBS.

## Alcance
- Migración de `hidden_lines` y su adaptador.

## Fuera
- Cancelaciones → FOKKO-0058B-SUBS.

## Features afectadas
- `internal/subscriptions` (`adapters/sqlite`)
- `internal/shared` (`sqlite/migrations`)

## Dudas abiertas
- `hidden_lines` no está en el modelo de datos de `ARCHITECTURE.md`: hay que sumarla cuando se actualice.

## Tareas
- [ ] [TASK-001 Tabla `hidden_lines`](TASK-001-HIDDEN-MIGRATION.md)
- [ ] [TASK-002 Adaptador de líneas ocultas](TASK-002-HIDDEN-LINES-REPO.md)
