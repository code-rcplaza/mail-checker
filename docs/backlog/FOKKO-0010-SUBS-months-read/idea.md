# FOKKO-0010-SUBS · Meses leídos de una cartola

## Problema
La subida debe responder "archivo y meses leídos", y Análisis lo muestra. Hoy `Statement` solo trae el número de cargos.

## Alcance
- Función pura que dice qué meses cubren los cargos de una cartola.
- `ImportStatement` devuelve esos meses en su resultado.

## Fuera
- Mostrarlos en pantalla → FOKKO-0100A-SUBS y FOKKO-0100B-SUBS.
- Persistirlos: se calculan, no se guardan.

## Features afectadas
- `internal/statements` (`domain`, `usecases`)

## Dudas abiertas
- El mes se cuenta en la hora de `America/Santiago`; hay que confirmar contra una cartola de cierre de mes.

## Tareas
- [x] [TASK-001A Meses distintos en orden](TASK-001A-MONTHS-DISTINCT.md)
- [x] [TASK-001B Mes según la hora de Santiago](TASK-001B-MONTH-BOUNDARY.md)
- [x] [TASK-002 `ImportStatement` devuelve los meses](TASK-002-IMPORT-RETURNS-MONTHS.md)
