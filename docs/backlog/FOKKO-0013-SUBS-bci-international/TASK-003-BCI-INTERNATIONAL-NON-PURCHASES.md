# FOKKO-0013-SUBS/TASK-003 · Ignorar lo que no es compra internacional
Feature: `internal/statements/adapters/bci`

Como usuario quiero que comisiones y traspasos del estado internacional no cuenten como cobros, para no ver suscripciones falsas.

## Criterios
- Dadas líneas de traspaso o comisión internacional, cuando se parsea, entonces no son cargos ni líneas malas.
- Dada una línea ilegible, cuando se parsea, entonces su número va a las líneas malas sin guardar el texto.

## Fuera
- Nada.

## Notas
- Ninguna.
