# 001-04 · Migraciones embebidas
Feature: `internal/shared/sqlite`

Como desarrollador quiero que las migraciones corran solas al arrancar para que cada
feature solo tenga que agregar su archivo `.sql`.

## Criterios
- Dada una base nueva y el directorio de migraciones embebido (sin tablas), cuando se
  llama `Migrate`, entonces termina sin error y existe la tabla de versiones de goose.
- Dada una base ya migrada, cuando se llama `Migrate` otra vez, entonces termina sin
  error y sin cambios (idempotente).
- Dado un FS de prueba con una migración SQL inválida, cuando se llama `Migrate`,
  entonces error que nombra el archivo de la migración.
- Dada una migración con dos sentencias donde la segunda falla, cuando se llama
  `Migrate`, entonces lo de la primera no queda aplicado (corre en transacción).

## Fuera
- Tablas de las features (fase 4).

## Notas
- `Migrate` recibe un `fs.FS` para poder testear con migraciones falsas; en
  producción se le pasa el embebido.
