# FOKKO-0011A-SUBS/TASK-001 · Reconocer una cartola BCI
Feature: `internal/statements/adapters/bci`

Como usuario quiero subir la cartola sin elegir el banco, para que la app lo reconozca sola.

## Criterios
- Dado el texto de una cartola BCI (nacional o internacional), cuando se llama `Detect`, entonces devuelve `true`.
- Dado el texto de otro documento, cuando se llama `Detect`, entonces devuelve `false`.

## Fuera
- Nada.

## Notas
- Fixtures reales anonimizados.
