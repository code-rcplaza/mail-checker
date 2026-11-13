# FOKKO-0105-SUBS/TASK-001 · `SubscriptionList`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver todas mis suscripciones juntas, para tener el panorama.

## Criterios
- Dadas cinco suscripciones, cuando se renderiza, entonces hay cinco filas ordenadas por próximo cobro.
- Dado un toque en una fila, cuando se hace, entonces lleva a su Detalle.
- Dado "Agregar", cuando se toca, entonces lleva a Agregar.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
