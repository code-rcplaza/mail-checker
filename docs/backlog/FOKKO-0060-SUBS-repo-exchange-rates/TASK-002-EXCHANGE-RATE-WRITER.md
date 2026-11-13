# FOKKO-0060-SUBS/TASK-002 · Guardar el tipo de cambio
Feature: `internal/statements/adapters/sqlite`

Como usuario quiero que subir una cartola guarde su tipo de cambio, para usarlo después.

## Criterios
- Dado un tipo de cambio con fecha, cuando se guarda dos veces, entonces queda uno (idempotente).
- Dados dos usuarios, cuando cada uno guarda el suyo, entonces no se mezclan.

## Fuera
- Nada.

## Notas
- Ninguna.
