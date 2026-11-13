# FOKKO-0098-SUBS/TASK-003 · Cerrar sesión
Feature: `frontend/src/features/auth/application`

Como usuario quiero cerrar mi sesión, para que nadie use mi cuenta.

## Criterios
- Dado "Cerrar sesión", cuando se toca, entonces llama a `POST /api/auth/logout`.
- Dado el cierre, cuando termina, entonces se vacía la caché de TanStack Query y se navega a Login.
- Dado un fallo de red al cerrar, cuando ocurre, entonces igual se limpia lo local y se avisa.

## Fuera
- Nada.

## Notas
- Caso de uso: se prueba el hook con el puerto doble (FOKKO-0085D-SUBS); no toca componentes.
