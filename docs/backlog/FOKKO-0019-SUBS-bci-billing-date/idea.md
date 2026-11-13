# FOKKO-0019-SUBS · Parser BCI: fecha de facturación

## Problema
BCI imprime la fecha en `FECHA ESTADO DE CUENTA` y el período en `PERÍODO FACTURADO`; el día de cierre cambia de un mes a otro, así que no se puede suponer.

## Depende de
- FOKKO-0011A-SUBS y FOKKO-0011B-SUBS y FOKKO-0018A-SUBS y FOKKO-0018B-SUBS.

## Alcance
- `BillingDate` del parser de BCI.

## Fuera
- Pedir la fecha al usuario cuando falta → FOKKO-0018A-SUBS y FOKKO-0018B-SUBS y FOKKO-0100A-SUBS y FOKKO-0100B-SUBS.

## Features afectadas
- `internal/statements` (`adapters/bci`)

## Dudas abiertas
- Formato exacto de las fechas en el texto: se confirma con las cartolas reales anonimizadas.

## Tareas
- [x] [TASK-001 Leer la fecha del estado de cuenta](TASK-001-BCI-READ-STATEMENT-DATE.md)
- [x] [TASK-002 Sin fecha en el archivo](TASK-002-BCI-BILLING-DATE-MISSING.md)
