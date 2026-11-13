# FOKKO-0061-SUBS/TASK-001 · Borrar lo vencido
Feature: `internal/auth/adapters/sqlite`

Como usuario quiero que la base no crezca sin parar, para que la app siga rápida.

## Criterios
- Dados códigos, sesiones y envíos vencidos y vigentes, cuando se limpia, entonces se borran solo los vencidos.
- Dado nada vencido, cuando se limpia, entonces no falla y no borra nada.
- Dados envíos de hace menos de una hora, cuando se limpia, entonces se conservan (cuentan para el límite).

## Fuera
- Nada.

## Notas
- Ninguna.
