# FOKKO-0106B-SUBS/TASK-001B · Accesibilidad del diálogo
Feature: `frontend/src/features/subscriptions/presentation/components`



## Criterios
- Dado el diálogo abierto, cuando se presiona Esc o se toca fuera, entonces se cierra sin hacer nada y el foco vuelve al botón.
- Dado el diálogo abierto, entonces el scroll de fondo queda bloqueado.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
