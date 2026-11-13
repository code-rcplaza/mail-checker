# FOKKO-0106B-SUBS/TASK-001A · Pedir confirmación
Feature: `frontend/src/features/subscriptions/presentation/components`



## Criterios
- Dado "Dejar de seguir", cuando se toca, entonces pide confirmación en un diálogo.
- Dada la confirmación, cuando se acepta, entonces llama a `DELETE /api/subscriptions/line` y vuelve a Subs sin esa línea.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
