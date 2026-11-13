# FOKKO-0098-SUBS · Protección de rutas

## Problema
Las rutas con sesión deben proteger el contenido; Login y Registro no deben mostrarse a quien ya entró.

## Depende de
- FOKKO-0068A-SUBS y FOKKO-0068B-SUBS y FOKKO-0085A-SUBS y FOKKO-0085B-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `useSession` (basado en `GET /api/me`), guardas de ruta y cierre de sesión.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/auth/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/auth`, `routes`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 `useSession`](TASK-001-USE-SESSION.md)
- [ ] [TASK-002 Guardas de ruta](TASK-002-ROUTE-GUARD.md)
- [ ] [TASK-003 Cerrar sesión](TASK-003-LOGOUT-ACTION.md)
- [ ] [TASK-004 Leer la sesión por HTTP](TASK-004-GATEWAY-READ-ME.md)
- [ ] [TASK-005 Cerrar sesión por HTTP](TASK-005-GATEWAY-LOGOUT.md)
