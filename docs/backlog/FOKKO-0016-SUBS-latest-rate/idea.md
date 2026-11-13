# FOKKO-0016-SUBS · Último tipo de cambio o total incompleto

## Problema
No todos los meses trae traspaso. Se guarda el último derivado con su fecha y se usa el más reciente; si no hay ninguno, el total avisa que está incompleto.

## Depende de
- FOKKO-0015-SUBS.

## Alcance
- Función pura que elige el tipo de cambio más reciente de una lista y dice si no hay ninguno.

## Fuera
- Guardarlo y leerlo → FOKKO-0060-SUBS y FOKKO-0118-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- ¿Se descarta un tipo de cambio demasiado viejo? Se propone que no (YAGNI): se muestra con su fecha.

## Tareas
- [ ] [TASK-001 Elegir el más reciente](TASK-001-PICK-LATEST-RATE.md)
- [ ] [TASK-002 Sin tipo de cambio](TASK-002-RATE-MISSING.md)
