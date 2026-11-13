# FOKKO-0051B-SUBS/TASK-001 · Tabla `code_sends`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que el límite de envíos se recuerde, para que un reinicio no lo reinicie.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces guarda correo y fecha de cada envío, ambos `NOT NULL`.

## Fuera
- Nada.

## Notas
- Ninguna.
