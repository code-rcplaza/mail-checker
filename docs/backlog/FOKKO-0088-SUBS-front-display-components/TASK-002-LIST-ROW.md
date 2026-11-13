# FOKKO-0088-SUBS/TASK-002 · `ListRow`
Feature: `frontend/src/components/ui`

Como usuario quiero filas fáciles de tocar, para entrar a cada detalle.

## Criterios
- Dada una fila con acción, cuando se toca o se activa con Enter, entonces se llama la acción.
- Dada una fila con título, subtítulo y valor, cuando se renderiza, entonces los tres se leen en ese orden.
- Dada una fila, entonces su alto táctil es de al menos 44 px.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
