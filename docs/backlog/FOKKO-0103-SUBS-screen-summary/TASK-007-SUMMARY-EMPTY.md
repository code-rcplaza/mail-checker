# FOKKO-0103-SUBS/TASK-007 · Resumen vacío
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero una guía si todavía no cargué nada, para empezar.

## Criterios
- Dado un resumen sin suscripciones, cuando se renderiza, entonces muestra el estado vacío con la acción "Subir cartola".
- Dada la acción, cuando se toca, entonces lleva a Elegir fuente.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
