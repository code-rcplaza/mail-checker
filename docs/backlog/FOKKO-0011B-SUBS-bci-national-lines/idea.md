# FOKKO-0011B-SUBS · Parser BCI: leer los cobros del estado nacional

## Problema
Sin un parser de BCI no se puede leer ninguna cartola real. El estado nacional trae las compras en pesos.

Esta idea cubre solo: Parser BCI: leer los cobros del estado nacional. El resto del trabajo original está en las ideas hermanas `FOKKO-0011*-SUBS`.

## Depende de
- FOKKO-0009-SUBS (el texto sale del extractor).

## Alcance
- Leer los cobros del estado nacional
- Ignorar lo que no es compra
- Líneas que no se pudieron leer

## Fuera
- Estado internacional → FOKKO-0013-SUBS.
- Fecha de facturación → FOKKO-0019-SUBS.
- Registro del banco en `main.go` → FOKKO-0082A-SUBS a FOKKO-0082D-SUBS.

## Features afectadas
- `internal/statements` (`adapters/bci`)

## Dudas abiertas
- Hay dos cartolas reales (sep y oct 2026). Se anonimizan antes de crear los fixtures, no después.
- ¿El parser deja el lugar pegado al descriptor (`HOSTGATOR LAS CONDES`) tal cual? FOKKO-0006-SUBS/TASK-002 reconoce por prefijo, así que se propone dejarlo.
- Compras en cuotas: no hay ejemplo real; hace falta uno para decidir si cada cuota es un cargo.

## Tareas
- [x] [TASK-001 Leer los cobros del estado nacional](TASK-001-BCI-NATIONAL-LINES.md)
- [x] [TASK-002 Ignorar lo que no es compra](TASK-002-BCI-SKIP-NON-PURCHASES.md)
- [x] [TASK-003 Líneas que no se pudieron leer](TASK-003-BCI-BAD-LINES.md)
