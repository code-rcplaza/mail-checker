# FOKKO-0031-SUBS/TASK-003 · Crear una suscripción manual
Feature: `internal/subscriptions/usecases`

Como usuario quiero agregar una suscripción a mano, para ver mi gasto completo.

## Criterios
- Dados datos válidos, cuando se crea, entonces se guarda con un `PublicID` generado por el caso de uso y se devuelve.
- Dados datos inválidos, cuando se crea, entonces no se guarda nada.
- Dados dos usuarios, cuando uno crea una, entonces el otro no la ve.

## Fuera
- Nada.

## Notas
- Ninguna.
