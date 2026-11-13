# FOKKO-0071-SUBS · `GET /api/banks`

## Problema
El selector de banco de la pantalla de subida necesita la lista de bancos soportados (el diseño no lo trae).

## Depende de
- FOKKO-0064-SUBS.

## Alcance
- `GET /api/banks` con los bancos registrados.

## Fuera
- Nada.

## Features afectadas
- `internal/statements` (`adapters/http`)

## Dudas abiertas
- El nombre visible del banco lo pone el front con i18n a partir del id (`bci`). ¿Se acepta?

## Tareas
- [ ] [TASK-001 `GET /api/banks`](TASK-001-BANKS-HANDLER.md)
