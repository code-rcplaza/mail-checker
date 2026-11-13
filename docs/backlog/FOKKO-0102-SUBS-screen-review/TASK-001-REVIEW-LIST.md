# FOKKO-0102-SUBS/TASK-001 · Lista de posibles
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver cada posible con su evidencia, para decidir con información.

## Criterios
- Dadas tres posibles, cuando se renderiza, entonces hay tres tarjetas con nombre, monto y "Lo vimos en tu archivo: …".
- Dado un monto en USD, cuando se renderiza, entonces muestra el monto original y su equivalente aproximado.
- Dada la lista, cuando se renderiza, entonces cada estado (pendiente, agregada, descartada) tiene su aspecto.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
