# FOKKO-0052-SUBS/TASK-001 · Tabla `sessions`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero seguir dentro después de que la app se reinicie, para no pedir otro código.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces tiene `user_id` con `ON DELETE CASCADE`, `token_hash` único y `expires_at`.
- Dado un usuario borrado, cuando se consulta `sessions`, entonces sus sesiones desaparecieron.

## Fuera
- Nada.

## Notas
- Ninguna.
