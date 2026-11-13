# FOKKO-0053C-SUBS/TASK-001 · Leer los cargos del usuario
Feature: `internal/subscriptions/adapters/sqlite`

Como usuario quiero que mi resumen use mis cargos guardados, para ver mis suscripciones.

## Criterios
- Dados cargos de dos usuarios, cuando `ListSince` pide los de uno, entonces solo devuelve los suyos.
- Dados cargos de distintas fechas, cuando se pide desde una fecha, entonces devuelve solo los posteriores o iguales.
- Dado un cargo guardado, cuando se lee, entonces todos sus campos coinciden, incluidos fecha, monto y últimos 4 dígitos.

## Fuera
- Nada.

## Notas
- Ninguna.
