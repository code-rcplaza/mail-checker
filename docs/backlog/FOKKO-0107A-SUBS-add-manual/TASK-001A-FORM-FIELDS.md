# FOKKO-0107A-SUBS/TASK-001A · Errores de los campos del formulario
Feature: `frontend/src/features/subscriptions/presentation/components`



## Criterios
- Dado un nombre vacío, cuando se valida, entonces hay error de nombre requerido.
- Dado un monto cero, negativo o no numérico, cuando se valida, entonces hay error de monto.
- Dado un día de cobro fuera de 1 a 31, cuando se valida, entonces hay error de día.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
