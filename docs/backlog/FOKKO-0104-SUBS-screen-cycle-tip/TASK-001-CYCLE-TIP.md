# FOKKO-0104-SUBS/TASK-001 · `CycleTip`
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero un consejo útil sobre mis fechas de cobro, para ordenar mis pagos.

## Criterios
- Dadas suscripciones con cobros concentrados en pocos días, cuando se renderiza, entonces muestra el tip con los días concentrados.
- Dado que no hay concentración, cuando se renderiza, entonces no se muestra.
- Dado el texto del tip, cuando se revisa, entonces no menciona cuentas bancarias ni acceso al banco.

## Fuera
- Nada.

## Notas
- Falta el texto aprobado por el usuario.
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
