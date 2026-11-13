# FOKKO-0025-SUBS/TASK-001 · Contar los restantes del mes
Feature: `internal/subscriptions/domain`

Como usuario quiero saber cuántos cobros me faltan este mes, para planificar.

## Criterios
- Dadas suscripciones con próximo cobro los días 20, 25 y 3 del mes siguiente, hoy día 10, cuando se cuenta, entonces son 2.
- Dado un próximo cobro hoy mismo, cuando se cuenta, entonces se documenta y prueba si cuenta como restante.
- Dada una suscripción ya cobrada este mes cuyo próximo cobro cae el mes siguiente, cuando se cuenta, entonces no entra.

## Fuera
- Nada.

## Notas
- Ninguna.
