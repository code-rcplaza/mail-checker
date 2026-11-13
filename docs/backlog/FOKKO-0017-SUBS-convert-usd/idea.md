# FOKKO-0017-SUBS · Convertir un cobro en US$ a CLP

## Problema
El total del mes va en CLP: un cobro en USD se muestra con su monto original y el equivalente en pesos con "≈".

## Depende de
- FOKKO-0015-SUBS.

## Alcance
- Función pura que convierte un monto en USD a CLP con un tipo de cambio.

## Fuera
- Armar el total → FOKKO-0029-SUBS.
- Mostrar el monto original junto al convertido → FOKKO-0030-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 Convertir dólares a pesos](TASK-001-USD-TO-CLP.md)
- [x] [TASK-002 Un monto en pesos no se toca](TASK-002-CLP-UNCHANGED.md)
