# FOKKO-0105-SUBS/TASK-002A · Datos básicos de la fila
Feature: `frontend/src/features/subscriptions/presentation/components`



## Criterios
- Dada una fila, cuando se renderiza, entonces muestra avatar, nombre (o alias), monto y próximo cobro.
- Dada cada fila, entonces muestra su banco y tarjeta (ej. "BCI · crédito").

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
