# FOKKO-0051A-SUBS/TASK-001 · Tabla `login_codes`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mi código sobreviva a un reinicio, para no quedar bloqueado.

## Criterios
- Dada la migración, cuando se aplica y se revierte en una base temporal, entonces funciona.
- Dada la tabla, entonces guarda hash, vencimiento e intentos por usuario, con un solo código vigente por usuario y borrado en cascada.
- Dado un segundo código para el mismo usuario insertado a mano, cuando se inserta, entonces la base lo rechaza.

## Fuera
- Nada.

## Notas
- Ninguna.
