# FOKKO-0087-SUBS/TASK-003 · `Checkbox`
Feature: `frontend/src/components/ui`

Como usuario quiero aceptar los términos con una casilla clara, para registrarme.

## Criterios
- Dada la casilla, cuando se toca la etiqueta, entonces se marca.
- Dada la casilla con teclado, cuando se presiona espacio, entonces cambia.
- Dado un error, cuando se renderiza, entonces se asocia a la casilla.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
