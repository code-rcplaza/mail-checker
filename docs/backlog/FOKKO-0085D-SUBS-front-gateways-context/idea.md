# FOKKO-0085D-SUBS · Puertos del front por Context

## Problema
Los casos de uso del front (`application`) dependen de puertos (gateways) y no de `fetch`. `main.tsx` entrega los reales por un Context y las pruebas entregan dobles en memoria, igual que en el back.

## Depende de
- FOKKO-0085A-SUBS y FOKKO-0085C-SUBS.

## Alcance
- `GatewaysProvider`, `useGateways` y el ayudante de pruebas `renderWithGateways`.

## Fuera
- Los gateways concretos de cada feature → la idea de cada pantalla.

## Features afectadas
- `frontend`

## Dudas abiertas
- La composición completa solo se puede probar cuando existan gateways (primeras pantallas).

## Tareas
- [ ] [TASK-001 Entregar los puertos con un Context](TASK-001-GATEWAYS-PROVIDER.md)
- [ ] [TASK-002 Fallar sin proveedor](TASK-002-GATEWAYS-MISSING.md)
- [ ] [TASK-003 Ayudante de pruebas con dobles](TASK-003-RENDER-WITH-GATEWAYS.md)
- [ ] [TASK-004 Armar los puertos reales en `main.tsx`](TASK-004-MAIN-COMPOSITION.md)
