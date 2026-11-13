# FOKKO-0030-SUBS · Monto original junto al equivalente ("≈")

## Problema
Un cobro en USD se muestra "USD 20 · ≈ $19.000": el monto original y su equivalente aproximado en pesos (decisión 19).

## Depende de
- FOKKO-0017-SUBS y FOKKO-0029-SUBS.

## Alcance
- Cada línea trae su monto original, su moneda y el equivalente en CLP con la marca de aproximado.

## Fuera
- Formato en pantalla → FOKKO-0105-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001A Mostrar un monto en dólares con su equivalente](TASK-001A-LINE-AMOUNT-USD.md)
- [x] [TASK-001B Mostrar un monto en pesos](TASK-001B-LINE-AMOUNT-CLP.md)
