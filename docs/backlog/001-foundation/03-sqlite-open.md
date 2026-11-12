# 001-03 · Abrir SQLite con PRAGMAs
Feature: `internal/shared/sqlite`

Como desarrollador quiero una sola forma de abrir la base, ya configurada, para que
ningún repo tenga que preocuparse de durabilidad ni de llaves foráneas.

## Criterios
- Dada una ruta en un directorio temporal, cuando se llama `Open`, entonces
  `journal_mode = wal`, `foreign_keys = 1`, `busy_timeout = 5000` y
  `synchronous = 2` (FULL).
- Dada la base abierta, cuando se revisa `db.Stats().MaxOpenConnections`, entonces
  es `1` (un solo escritor).
- Dada una ruta cuyo directorio padre no existe, cuando se llama `Open`, entonces se
  crea el directorio y la base abre bien.
- Dada una ruta no escribible, cuando se llama `Open`, entonces error que incluye
  la ruta.

## Fuera
- Migraciones (001-04), transacciones (001-05).

## Notas
- Los PRAGMAs se aplican por conexión (vía DSN de `modernc.org/sqlite`), no con un
  `Exec` suelto que solo afectaría a una conexión.
