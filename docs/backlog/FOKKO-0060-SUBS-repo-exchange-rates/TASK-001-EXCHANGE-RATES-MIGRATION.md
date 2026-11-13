# FOKKO-0060-SUBS/TASK-001 · Tabla `exchange_rates`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mi tipo de cambio quede guardado, para convertir mis dólares.

## Criterios
- Dada la migración, cuando se aplica y se revierte, entonces funciona.
- Dada la tabla, entonces guarda usuario, tipo de cambio, fecha y `UNIQUE(user_id, date)`, con borrado en cascada.
- Dado un tipo de cambio cero o negativo, cuando se inserta, entonces la base lo rechaza.

## Fuera
- Nada.

## Notas
- Ninguna.
