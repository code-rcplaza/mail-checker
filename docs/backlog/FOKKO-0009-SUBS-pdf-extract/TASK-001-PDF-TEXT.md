# FOKKO-0009-SUBS/TASK-001 · Extraer el texto de un PDF
Feature: `internal/statements/adapters/extract`

Como usuario quiero que la app lea el texto de mi cartola en PDF, para que se puedan detectar mis cobros.

## Criterios
- Dado un PDF con texto, cuando se extrae, entonces devuelve el texto con las líneas en el orden de la página y el tipo `pdf`.
- Dado un PDF de varias páginas, cuando se extrae, entonces el texto de las páginas sale en orden.

## Fuera
- Nada.

## Notas
- Fixture: PDF real anonimizado antes de crearlo.
