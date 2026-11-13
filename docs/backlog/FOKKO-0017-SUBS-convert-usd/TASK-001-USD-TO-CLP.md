# FOKKO-0017-SUBS/TASK-001 · Convertir dólares a pesos
Feature: `internal/subscriptions/domain`

Como usuario quiero ver mis cobros en dólares en pesos, para compararlos con el resto.

## Criterios
- Dado US$ 20,00 y un tipo de cambio de 950, cuando se convierte, entonces son CLP 19.000 y el resultado es aproximado.
- Dado un monto con decimales al convertir, cuando se convierte, entonces redondea half-up a centésimas.
- Dado US$ 0, cuando se convierte, entonces son CLP 0.

## Fuera
- Nada.

## Notas
- Ninguna.
