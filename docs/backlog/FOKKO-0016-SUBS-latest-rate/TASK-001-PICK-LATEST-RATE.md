# FOKKO-0016-SUBS/TASK-001 · Elegir el más reciente
Feature: `internal/subscriptions/domain`

Como usuario quiero que se use el tipo de cambio más nuevo, para que el total sea lo más cercano a lo que pago.

## Criterios
- Dados tres tipos de cambio con fechas distintas, cuando se elige, entonces devuelve el de fecha más reciente.
- Dados dos con la misma fecha, cuando se elige, entonces el resultado es determinista (se documenta el desempate).
- Dada una lista que no se modifica, cuando se elige, entonces la lista original queda intacta.

## Fuera
- Nada.

## Notas
- Ninguna.
