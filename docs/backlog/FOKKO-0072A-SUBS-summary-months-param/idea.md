# FOKKO-0072A-SUBS · `GET /api/summary`: parámetro `months`

## Problema
`GET /api/summary` es la lectura más grande del MVP y se parte por sección. Esta idea es la base: la ruta, la sesión y el parámetro `months`.

## Depende de
- FOKKO-0021-SUBS y FOKKO-0064-SUBS.

## Alcance
- Handler con sesión y validación de `months`; devuelve el resumen vacío si no hay datos.

## Fuera
- Cada sección del resumen → FOKKO-0072B-SUBS a FOKKO-0072H-SUBS.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Convenciones de la respuesta: montos en centésimas (entero) con su moneda, fechas en unix segundos y `lineKey` como cadena opaca.

## Tareas
- [x] [TASK-001 Parámetro `months`](TASK-001-SUMMARY-QUERY-PARAMS.md)
- [x] [TASK-002 Usuario sin datos](TASK-002-SUMMARY-EMPTY-USER.md)
