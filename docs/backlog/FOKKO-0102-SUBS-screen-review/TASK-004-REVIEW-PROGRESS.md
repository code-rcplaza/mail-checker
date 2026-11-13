# FOKKO-0102-SUBS/TASK-004 · Progreso de la revisión
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver cuántas me faltan, para saber cuándo termino.

## Criterios
- Dadas 5 posibles con 2 resueltas, cuando se renderiza, entonces el anillo muestra 2 de 5.
- Dadas todas resueltas, cuando se renderiza, entonces aparece la acción "Ver mi resumen".
- Dado ningún posible, cuando se abre, entonces muestra el estado vacío (Revisión vacía) y lleva al Resumen.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
