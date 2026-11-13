# FOKKO-0090-SUBS/TASK-003 · `SegmentedControl`
Feature: `frontend/src/components/ui`

Como usuario quiero elegir una opción entre pocas (mensual o anual), para decidir rápido.

## Criterios
- Dadas tres opciones, cuando se elige una, entonces se llama el manejador con su valor.
- Dado el teclado, cuando se usan las flechas, entonces se mueve entre opciones (rol de grupo de radio).
- Dada una opción seleccionada, cuando se renderiza, entonces queda marcada para lectores de pantalla.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
