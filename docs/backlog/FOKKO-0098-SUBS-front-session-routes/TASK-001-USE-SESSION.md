# FOKKO-0098-SUBS/TASK-001 · `useSession`
Feature: `frontend/src/features/auth/application`

Como usuario quiero que la app sepa si estoy dentro, para mostrarme lo mío.

## Criterios
- Dado un 200 de `GET /api/me`, cuando se lee, entonces la sesión trae nombre, correo e iniciales.
- Dado un 401, cuando se lee, entonces la sesión es "sin sesión" y no es un error visible.
- Dado el estado de carga, cuando se lee, entonces se distingue de "sin sesión" (no parpadea Login).

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
