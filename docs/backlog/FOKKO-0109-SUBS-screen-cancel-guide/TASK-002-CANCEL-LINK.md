# FOKKO-0109-SUBS/TASK-002 · Enlace al sitio del servicio
Feature: `frontend/src/features/subscriptions/presentation/components`

Como usuario quiero ir directo al sitio, para cancelar en un paso.

## Criterios
- Dado `{url}` con `https`, cuando se renderiza, entonces el botón abre en otra pestaña con `rel="noopener noreferrer"`.
- Dado `url` vacío, cuando se renderiza, entonces no hay botón y se muestra un consejo genérico.
- Dado un `url` que no es `https`, cuando se renderiza, entonces no se muestra como enlace.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
