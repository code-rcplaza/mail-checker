# FOKKO-0014-SUBS · Parser BCI: par de traspasos

## Problema
El tipo de cambio sale del par de traspasos: `TRASPASO DEUDA INTERNACIONAL` (CLP, estado nacional) y `TRASPASO DEUDA NACIONAL` (US$, estado internacional). El parser tiene que entregarlos.

## Depende de
- FOKKO-0011A-SUBS y FOKKO-0011B-SUBS y FOKKO-0013-SUBS.

## Alcance
- El parser de BCI entrega, si existe, el par de montos y su fecha.

## Fuera
- Calcular el tipo de cambio → FOKKO-0015-SUBS.
- Guardarlo → FOKKO-0118-SUBS.

## Features afectadas
- `internal/statements` (`ports`, `adapters/bci`)

## Dudas abiertas
- Cómo viaja el par: un método nuevo en `BankStatementParser` o ampliar el resultado de `Parse`. Se propone un método propio, para no cambiar la firma de `Parse`.
- No todos los meses trae traspaso (en uno el monto era US$ 0,00).

## Tareas
- [x] [TASK-001 Entregar el par desde el puerto](TASK-001-TRANSFER-PAIR-PORT.md)
- [x] [TASK-002 Leer el par de traspasos en BCI](TASK-002-BCI-READ-TRANSFER-PAIR.md)
