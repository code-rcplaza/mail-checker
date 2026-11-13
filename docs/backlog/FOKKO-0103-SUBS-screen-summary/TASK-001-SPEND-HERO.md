# FOKKO-0103-SUBS/TASK-001 · `SpendHero`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ver cuánto gasto al mes de un vistazo, para tener una referencia.

## Criterios
- Dado el total en CLP, cuando se renderiza, entonces el anillo y el número animado muestran el monto formateado en es-CL.
- Dado un total aproximado (hay dólares convertidos), cuando se renderiza, entonces lleva "≈".
- Dado un total incompleto (faltó el tipo de cambio), cuando se renderiza, entonces avisa que no incluye los dólares.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
