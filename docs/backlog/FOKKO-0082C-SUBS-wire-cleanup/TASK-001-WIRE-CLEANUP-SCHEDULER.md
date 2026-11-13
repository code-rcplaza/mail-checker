# FOKKO-0082C-SUBS/TASK-001 · Limpieza periódica
Feature: `cmd/api`

Como usuario quiero que la base no se llene de códigos viejos, para que la app siga ligera.

## Criterios
- Dado el arranque, cuando pasa el intervalo, entonces se ejecuta la limpieza de lo vencido.
- Dado el cierre de la app, cuando se cancela el contexto, entonces la tarea termina limpiamente.
- Dado un error en la limpieza, cuando ocurre, entonces se loguea sin datos personales y la tarea sigue.

## Fuera
- Nada.

## Notas
- Ninguna.
