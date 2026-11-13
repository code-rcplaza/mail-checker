# FOKKO-0059-SUBS/TASK-002 · Adaptador de alzas descartadas
Feature: `internal/subscriptions/adapters/sqlite`

Como usuario quiero que solo se oculte el alza que ya vi, para enterarme de la próxima.

## Criterios
- Dado un descarte para el monto 10.500, cuando se lista, entonces aparece solo ese monto para esa línea.
- Dado otro descarte para 11.000 en la misma línea, cuando se lista, entonces aparecen los dos.
- Dados descartes de dos usuarios, cuando uno lista, entonces ve solo los suyos.

## Fuera
- Nada.

## Notas
- Ninguna.
