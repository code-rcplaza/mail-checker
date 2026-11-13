# FOKKO-0100A-SUBS/TASK-001 · Cargando
Feature: `frontend/src/features/statements/presentation/components`

Como usuario quiero ver que mi cartola se está leyendo, para no pensar que se colgó.

## Criterios
- Dada la subida pendiente, cuando se renderiza, entonces muestra el anillo de carga y "Paso 2 de 3".
- Dada la subida pendiente, cuando se renderiza, entonces no se puede enviar otra.
- Dada la subida terminada, cuando responde, entonces navega a Análisis.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
