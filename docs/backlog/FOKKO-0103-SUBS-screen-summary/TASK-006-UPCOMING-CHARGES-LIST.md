# FOKKO-0103-SUBS/TASK-006 · `UpcomingChargesList`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver qué cobros vienen y cuándo, para llegar con saldo.

## Criterios
- Dados varios próximos cobros, cuando se renderiza, entonces salen ordenados por fecha.
- Dado un toque en una fila, cuando se hace, entonces lleva al Detalle de esa línea.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
