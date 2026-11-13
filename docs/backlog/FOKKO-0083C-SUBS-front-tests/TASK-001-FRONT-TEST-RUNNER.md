# FOKKO-0083C-SUBS/TASK-001 · Pruebas de componentes
Feature: `frontend`

Como usuario quiero que cada pantalla tenga pruebas, para cambiar sin romper.

## Criterios
- Dado un componente de ejemplo, cuando se corre `pnpm test`, entonces se renderiza y se comprueba su texto.
- Dado un test fallido, cuando se corre `pnpm test`, entonces el comando termina con error.
- Dado el entorno de pruebas, entonces simula `fetch` sin red real.

## Fuera
- Nada.

## Notas
- El componente de ejemplo se borra al terminar la tarea.
