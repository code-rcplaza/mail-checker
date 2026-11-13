# FOKKO-0012-SUBS · Decidir si un cobro internacional es en pesos

## Problema
La cartola internacional no trae el código de la moneda de origen: trae el monto de origen y el monto en US$. Si el origen es muchísimo mayor que el US$, el cobro fue en pesos y hay que usar el origen tal cual (decisión 6).

## Alcance
- Función pura con una constante con nombre (razón de cientos) que decide si un cobro internacional es en pesos.

## Fuera
- Leer las columnas de la cartola → FOKKO-0013-SUBS.

## Features afectadas
- `internal/statements` (`domain`)

## Dudas abiertas
- Valor de la constante: se propone 100. Hay que validarlo con las cartolas reales antes de fijarlo.

## Tareas
- [x] [TASK-001A Decidir por la razón origen sobre US$](TASK-001A-PESOS-RATIO-CASES.md)
- [x] [TASK-001B Límite exacto de la razón](TASK-001B-PESOS-RATIO-BOUNDARY.md)
- [x] [TASK-002 US$ en cero](TASK-002-ZERO-USD-AMOUNT.md)
