# FOKKO-0060-SUBS/TASK-003 · Leer el más reciente
Feature: `internal/subscriptions/adapters/sqlite`

Como usuario quiero que mi resumen use el último tipo de cambio, para ver el total correcto.

## Criterios
- Dados tres tipos de cambio de fechas distintas, cuando se lee el último, entonces devuelve el de la fecha más reciente.
- Dado ninguno, cuando se lee, entonces devuelve que no hay.
- Dados tipos de cambio de dos usuarios, cuando uno lee, entonces ve solo los suyos.

## Fuera
- Nada.

## Notas
- Ninguna.
