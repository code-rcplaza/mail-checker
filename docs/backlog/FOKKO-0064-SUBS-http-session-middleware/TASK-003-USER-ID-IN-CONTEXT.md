# FOKKO-0064-SUBS/TASK-003 · El `userID` en el contexto
Feature: `internal/auth/adapters/http`

Como usuario quiero que cada handler sepa quién soy, para que solo toque mis datos.

## Criterios
- Dado un `userID` guardado por el middleware, cuando un handler lo lee, entonces lo obtiene.
- Dado un handler que lo lee sin que haya pasado por el middleware, cuando lo lee, entonces recibe un error y no un `userID` cero.

## Fuera
- Nada.

## Notas
- Ninguna.
