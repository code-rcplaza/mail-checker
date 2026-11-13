# FOKKO-0093-SUBS/TASK-002 · Corte de escritorio
Feature: `frontend/src/components/layout`

Como usuario quiero que la app se adapte a mi pantalla, para verla bien en cualquier dispositivo.

## Criterios
- Dado el shell, cuando se renderiza, entonces incluye `BottomNav` visible por debajo de `lg:` y `SideNav` visible desde `lg:`.
- Dado el ancho de 768 a 1023 px, entonces usa el layout de escritorio con la barra inferior.
- Dadas las pantallas de acceso, cuando se renderizan, entonces no muestran barra lateral ni inferior.

## Fuera
- Nada.

## Notas
- Componente presentacional: solo props y `t()`; sin `useQuery`, router ni `fetch` (FOKKO-0085C-SUBS).
