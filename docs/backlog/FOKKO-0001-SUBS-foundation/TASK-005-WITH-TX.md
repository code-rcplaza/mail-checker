# FOKKO-0001-SUBS/TASK-005 · Transacciones con `WithTx`
Feature: `internal/shared/sqlite`

Como desarrollador quiero un helper de transacción para que cada repo escriba
varias filas de forma atómica sin repetir el manejo de commit y rollback.

## Criterios
- Dada una función que inserta una fila y devuelve `nil`, cuando corre en `WithTx`,
  entonces la fila queda guardada.
- Dada una función que inserta una fila y devuelve un error, cuando corre en
  `WithTx`, entonces la fila no queda y `WithTx` devuelve ese mismo error (`errors.Is`).
- Dada una función que inserta una fila y hace `panic`, cuando corre en `WithTx`,
  entonces se hace rollback y el `panic` se propaga.
- Dado un contexto ya cancelado, cuando se llama `WithTx`, entonces devuelve error y
  no ejecuta la función.

## Fuera
- Repositorios reales (fase 4).

## Notas
- Los tests usan una tabla creada solo para el test.
