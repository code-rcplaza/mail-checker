# FOKKO-0088-SUBS · `Card`, `Badge`, `ListRow`, `ServiceAvatar`

## Problema
Las listas y tarjetas del MVP comparten componentes de presentación.

## Depende de
- FOKKO-0084A-SUBS y FOKKO-0084B-SUBS.
- FOKKO-0085C-SUBS.

## Alcance
- `Card` (hero, alerta y tip), `Badge`, `ListRow` y `ServiceAvatar`.

## Fuera
- Nada.

## Features afectadas
- `frontend` (`components/ui`)

## Dudas abiertas
- Logos de servicios: se propone no usarlos (derechos de marca) y mostrar la inicial con un color estable por nombre.
  - Resuelto: sí (opción propuesta); el color sale de los tokens `series-1…6` (colores de servicio del handoff).

## Tareas
- [x] [TASK-001A Variantes de `Card`](TASK-001A-CARD-VARIANTS.md)
- [x] [TASK-001B Variantes de `Badge`](TASK-001B-BADGE-VARIANTS.md)
- [x] [TASK-002 `ListRow`](TASK-002-LIST-ROW.md)
- [x] [TASK-003 `ServiceAvatar`](TASK-003-SERVICE-AVATAR.md)
