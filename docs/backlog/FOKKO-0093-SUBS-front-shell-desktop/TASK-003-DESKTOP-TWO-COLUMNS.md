# FOKKO-0093-SUBS/TASK-003 · Contenido en dos columnas
Feature: `frontend/src/components/layout`

Como usuario quiero aprovechar la pantalla grande, para ver el resumen y el detalle a la vez.

## Criterios
- Dado el contenido con sesión en escritorio, cuando se renderiza, entonces mide como máximo 1200 px y tiene dos columnas con 48 px de separación.
- Dada la columna izquierda, entonces aloja el hero y la derecha el contenido de la pantalla.
- Dado un ancho móvil, cuando se renderiza, entonces las columnas se apilan.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
