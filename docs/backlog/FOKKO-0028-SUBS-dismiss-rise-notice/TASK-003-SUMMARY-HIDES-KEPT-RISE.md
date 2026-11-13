# FOKKO-0028-SUBS/TASK-003 · El resumen oculta el alza descartada
Feature: `internal/subscriptions/domain`

Como usuario quiero que un alza ya vista no vuelva a aparecer, para no ver ruido.

## Criterios
- Dada un alza descartada para el monto 10.500, cuando se arma el resumen, entonces no aparece.
- Dada la misma línea con un nuevo cobro de 11.000, cuando se arma el resumen, entonces aparece la nueva alza.

## Fuera
- Nada.

## Notas
- Ninguna.
