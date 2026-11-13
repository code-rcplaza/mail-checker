# FOKKO-0011B-SUBS/TASK-003 · Líneas que no se pudieron leer
Feature: `internal/statements/adapters/bci`

Como usuario quiero saber si alguna línea no se pudo leer, para sospechar de un archivo raro.

## Criterios
- Dada una línea con forma de compra pero con monto o fecha ilegible, cuando se parsea, entonces su número va a la lista de líneas malas.
- Dada una línea mala, cuando se parsea, entonces el texto de la línea no aparece en ningún resultado.

## Fuera
- Nada.

## Notas
- Ninguna.
