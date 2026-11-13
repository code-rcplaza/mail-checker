# FOKKO-0103-SUBS/TASK-005 · `PendingActionsList`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero saber qué me falta resolver, para dejar todo al día.

## Criterios
- Dadas alzas y posibles pendientes, cuando se renderiza, entonces cada una es una fila que lleva a Alerta o a Revisión.
- Dado ningún pendiente, cuando se renderiza, entonces no se muestra la sección.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
