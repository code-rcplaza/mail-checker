# FOKKO-0084B-SUBS · Preferencia de tema

## Problema
El diseño usa tokens de color y tema claro/oscuro. El tema vive en el navegador (sin endpoint).

Esta idea cubre solo: Preferencia de tema. El resto del trabajo original está en las ideas hermanas `FOKKO-0084*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS.

## Alcance
- Preferencia de tema

## Fuera
- La pantalla del selector → FOKKO-0110A-SUBS y FOKKO-0110B-SUBS.

## Features afectadas
- `frontend` (`lib/theme`, estilos)

## Dudas abiertas
- Los tokens salen del `tokens.json` del handoff Fokko Subs. Hace falta tener ese archivo a mano; no se inventan colores.

## Tareas
- [x] [TASK-001A Seguir el tema del sistema](TASK-001A-THEME-SYSTEM.md)
- [x] [TASK-001B Recordar la preferencia](TASK-001B-THEME-PERSIST.md)
