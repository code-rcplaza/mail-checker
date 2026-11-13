# FOKKO-0031-SUBS/TASK-004 · La manual aparece en el resumen
Feature: `internal/subscriptions/domain`

Como usuario quiero que mi suscripción manual cuente en el total, para ver el gasto verdadero.

## Criterios
- Dada una manual de CLP 5.000, cuando se arma el resumen, entonces aparece en suscripciones y suma al gasto mensual.
- Dada una manual en USD, cuando se arma el resumen, entonces se convierte como cualquier cobro en USD.
- Dada una manual, cuando se arma el resumen, entonces su próximo cobro es el de su día.

## Fuera
- Nada.

## Notas
- Ninguna.
