# FOKKO-0099-SUBS/TASK-001A · Elegir o soltar un PDF
Feature: `frontend/src/features/statements/presentation/components`



## Criterios
- Dado un archivo PDF, cuando se suelta o se elige, entonces queda seleccionado y se muestra su nombre.
- Dado el teclado, cuando se activa la zona con Enter o espacio, entonces abre el selector de archivos.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
