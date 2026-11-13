# FOKKO-0013-SUBS · Parser BCI: cobros del estado internacional (origen y US$)

## Problema
Las compras en el extranjero van en el estado internacional de la cartola BCI. Sin leerlo, faltan Claude, Netflix en USD y otros.

## Depende de
- FOKKO-0011A-SUBS y FOKKO-0011B-SUBS (`Detect` y la base del parser).
- FOKKO-0012-SUBS (decide pesos o US$).

## Alcance
- Leer las líneas del estado internacional: monto de origen y monto en US$.
- Guardar el cargo en CLP (origen) o en USD (US$) según FOKKO-0012-SUBS.

## Fuera
- Conversión a pesos para mostrar → FOKKO-0017-SUBS.
- Par de traspasos → FOKKO-0014-SUBS.

## Features afectadas
- `internal/statements` (`adapters/bci`)

## Dudas abiertas
- Un cobro en euros u otra moneda no se ha visto en una cartola real; su test es sintético hasta que aparezca uno.
- Fixtures reales anonimizados (hay dos cartolas con estado internacional).

## Tareas
- [x] [TASK-001 Leer las líneas del estado internacional](TASK-001-BCI-INTERNATIONAL-LINES.md)
- [x] [TASK-002 Guardar en pesos o en dólares](TASK-002-BCI-CURRENCY-BY-RATIO.md)
- [x] [TASK-003 Ignorar lo que no es compra internacional](TASK-003-BCI-INTERNATIONAL-NON-PURCHASES.md)
