# FOKKO-0112B-SUBS/TASK-001A · Subir la cartola
Feature: `cmd/api`



## Criterios
- Dada la cartola anonimizada, cuando se sube, entonces `POST /api/statements` responde 201 con los meses leídos.
- Dada la misma cartola, cuando se sube otra vez, entonces responde 409 `statements.already_imported`.

## Fuera
- Nada.

## Notas
- Ninguna.
