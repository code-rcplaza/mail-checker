# FOKKO-0118-SUBS/TASK-002 · Importar guarda el tipo de cambio
Feature: `internal/statements/usecases`

Como usuario quiero que subir mi cartola actualice el tipo de cambio, para que el total sea el de mi último mes.

## Criterios
- Dada una cartola con par de traspasos, cuando se importa, entonces se guarda el tipo de cambio con la fecha de facturación.
- Dada una cartola sin par, cuando se importa, entonces no se guarda nada y la importación sigue.
- Dado un fallo al guardar el tipo de cambio, cuando se importa, entonces se informa el error sin dejar la cartola a medias.

## Fuera
- Nada.

## Notas
- Ninguna.
