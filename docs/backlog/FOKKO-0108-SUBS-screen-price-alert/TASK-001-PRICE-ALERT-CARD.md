# FOKKO-0108-SUBS/TASK-001 · `PriceAlertCard`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver cuánto subió y cuánto me cuesta, para decidir.

## Criterios
- Dada una alza, cuando se renderiza, entonces muestra el servicio, el antes, el después, desde cuándo y el impacto mensual y anual.
- Dada una alza que es cambio de plan, cuando se renderiza, entonces lo indica.
- Dado un monto en USD, cuando se renderiza, entonces muestra el equivalente aproximado.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
