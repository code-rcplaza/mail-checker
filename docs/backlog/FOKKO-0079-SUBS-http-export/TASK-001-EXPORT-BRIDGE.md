# FOKKO-0079-SUBS/TASK-001 · Juntar los datos del usuario
Feature: `internal/auth/adapters/export`

Como usuario quiero llevarme mis datos, para tener una copia.

## Criterios
- Dado un usuario con datos en varias features, cuando se arma la copia, entonces trae solo los suyos.
- Dado un usuario sin datos, cuando se arma, entonces trae las listas vacías.
- Dado el puente, entonces es el único lugar que importa casos de uso de otras features.

## Fuera
- Nada.

## Notas
- Ninguna.
