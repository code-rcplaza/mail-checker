# FOKKO-0103-SUBS/TASK-004 · `NextChargeRow`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver mi próximo cobro, para estar preparado.

## Criterios
- Dado el próximo cobro de Netflix el día 5, cuando se renderiza, entonces muestra nombre, monto y fecha (es-CL, hora de Santiago).
- Dado que no hay próximos cobros, cuando se renderiza, entonces no se muestra la fila.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
