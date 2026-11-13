# FOKKO-0011B-SUBS/TASK-001 · Leer los cobros del estado nacional
Feature: `internal/statements/adapters/bci`

Como usuario quiero que mis compras en pesos se lean de la cartola, para ver mis suscripciones.

## Criterios
- Dada una línea de compra, cuando se parsea, entonces devuelve monto, descriptor, fecha, tarjeta y últimos 4 dígitos.
- Dado el monto `$5.990`, cuando se parsea, entonces vale `599000` centésimas en CLP.
- Dada una cartola con varias compras, cuando se parsea, entonces cada línea trae su número de línea del texto.

## Fuera
- Nada.

## Notas
- Un descriptor real con teléfono o ciudad se conserva; no se inventan descriptores.
