# FOKKO-0103-SUBS/TASK-003 · Cobros del mes
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver cuántos cobros llevo y cuántos me quedan, para planificar el mes.

## Criterios
- Dados 3 cobrados y 5 restantes, cuando se renderiza, entonces dice "Llevas 3 cobros este mes y te quedan 5".
- Dado 1 cobrado, cuando se renderiza, entonces usa el singular correcto.
- Dados 0 y 0, cuando se renderiza, entonces no muestra la línea.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
