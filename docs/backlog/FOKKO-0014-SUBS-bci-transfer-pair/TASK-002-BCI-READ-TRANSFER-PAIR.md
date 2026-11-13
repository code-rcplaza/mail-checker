# FOKKO-0014-SUBS/TASK-002 · Leer el par de traspasos en BCI
Feature: `internal/statements/adapters/bci`

Como usuario quiero que se lea el traspaso de mi cartola, para derivar el tipo de cambio.

## Criterios
- Dada una cartola con ambos traspasos, cuando se lee el par, entonces trae el CLP del nacional y el US$ del internacional.
- Dada una cartola sin traspasos, cuando se lee el par, entonces devuelve que no hay par y no es un error.
- Dado un traspaso con US$ 0,00, cuando se lee el par, entonces devuelve que no hay par.

## Fuera
- Nada.

## Notas
- Ninguna.
