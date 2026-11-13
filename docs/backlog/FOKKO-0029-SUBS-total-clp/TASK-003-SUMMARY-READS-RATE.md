# FOKKO-0029-SUBS/TASK-003 · `GetSummary` lee el tipo de cambio
Feature: `internal/subscriptions/usecases`

Como usuario quiero que el resumen use mi último tipo de cambio, para ver un total convertido.

## Criterios
- Dado un puerto `ExchangeRates` con un tipo de cambio guardado, cuando se arma el resumen, entonces el total lo usa.
- Dado un puerto sin tipos de cambio, cuando se arma el resumen, entonces el total sale incompleto.
- Dado un error del puerto, cuando se arma el resumen, entonces devuelve el error con contexto.

## Fuera
- Nada.

## Notas
- Ninguna.
