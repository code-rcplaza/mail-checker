# FOKKO-0036-SUBS · `Summary`: reparto por servicio

## Problema
El Resumen muestra el total del mes y un anillo con el reparto por servicio.

## Depende de
- FOKKO-0029-SUBS.

## Alcance
- Reparto del gasto mensual en CLP por servicio (padre), ordenado.

## Fuera
- Dibujar el anillo → FOKKO-0089A-SUBS y FOKKO-0089B-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Cuántos segmentos muestra el anillo: se propone los 5 mayores y el resto junto en "Otros".

## Tareas
- [ ] [TASK-001 Reparto por servicio](TASK-001-SPEND-SHARE-BY-PARENT.md)
- [ ] [TASK-002 Cinco mayores más "Otros"](TASK-002-SPEND-TOP-AND-OTHERS.md)
