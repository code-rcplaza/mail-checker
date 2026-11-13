# FOKKO-0106A-SUBS/TASK-003 · `ChargeHistory`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver mis últimos cobros, para verificar montos y fechas.

## Criterios
- Dados seis cobros, cuando se renderiza, entonces salen del más reciente al más antiguo con fecha y monto.
- Dado un cobro en USD, cuando se renderiza, entonces muestra el monto original y su equivalente.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
