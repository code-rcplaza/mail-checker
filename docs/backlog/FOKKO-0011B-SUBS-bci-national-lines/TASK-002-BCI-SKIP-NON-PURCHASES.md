# FOKKO-0011B-SUBS/TASK-002 · Ignorar lo que no es compra
Feature: `internal/statements/adapters/bci`

Como usuario quiero que pagos, intereses y traspasos no aparezcan como cobros, para no ver suscripciones falsas.

## Criterios
- Dadas líneas de pago, interés, comisión o traspaso de deuda, cuando se parsea, entonces no se devuelven como cargos.
- Dadas esas líneas, cuando se parsea, entonces tampoco se cuentan como líneas malas.

## Fuera
- Nada.

## Notas
- Ninguna.
