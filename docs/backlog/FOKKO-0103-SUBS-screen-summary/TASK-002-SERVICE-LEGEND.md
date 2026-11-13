# FOKKO-0103-SUBS/TASK-002 · `ServiceLegend`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero saber qué servicio pesa más, para decidir dónde recortar.

## Criterios
- Dado el reparto con 5 servicios y "Otros", cuando se renderiza, entonces la leyenda lista cada uno con su porcentaje y color del anillo.
- Dado un solo servicio, cuando se renderiza, entonces la leyenda tiene una fila.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
