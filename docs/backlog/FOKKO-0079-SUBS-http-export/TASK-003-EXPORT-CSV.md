# FOKKO-0079-SUBS/TASK-003 · Copia en CSV
Feature: `internal/auth/adapters/export`

Como usuario quiero ver mis suscripciones en una planilla, para revisarlas.

## Criterios
- Dadas mis suscripciones, cuando se serializan, entonces cada una es una fila con encabezado.
- Dado un campo que empieza con `=`, `+`, `-` o `@`, cuando se serializa, entonces se neutraliza para que no se ejecute como fórmula.
- Dado un campo con coma, comillas o salto de línea, cuando se serializa, entonces se escapa según CSV.

## Fuera
- Nada.

## Notas
- Ninguna.
