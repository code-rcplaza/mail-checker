# FOKKO-0100B-SUBS/TASK-002 · Sin suscripciones encontradas
Feature: `frontend/src/features/statements/presentation/components`

Como usuario quiero un mensaje claro si no se encontró ninguna suscripción, para saber qué hacer.

## Criterios
- Dado un resumen sin suscripciones ni posibles, cuando se renderiza Análisis, entonces muestra el estado vacío con la opción de subir otra cartola.
- Dado ese estado, cuando se toca "Subir otra", entonces vuelve a Elegir fuente.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
