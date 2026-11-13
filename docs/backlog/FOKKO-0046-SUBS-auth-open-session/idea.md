# FOKKO-0046-SUBS · Crear sesión de 30 días

## Problema
Verificar el código abre una sesión de 30 días. No hay contraseña: la sesión es un token opaco que viaja en una cookie.

## Depende de
- FOKKO-0043-SUBS.

## Alcance
- Dominio `Session` y puerto `Sessions` (se guarda el hash del token).
- Caso de uso que abre la sesión tras verificar el código.

## Fuera
- Validar y cerrar → FOKKO-0047A-SUBS y FOKKO-0047B-SUBS.
- La cookie → FOKKO-0065A-SUBS a FOKKO-0065C-SUBS.

## Features afectadas
- `internal/auth` (`domain`, `ports`, `usecases`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Generar el token](TASK-001A-SESSION-TOKEN.md)
- [ ] [TASK-001B Fijar el vencimiento](TASK-001B-SESSION-EXPIRY.md)
- [ ] [TASK-001C Guardar solo el hash](TASK-001C-SESSION-HASH.md)
- [ ] [TASK-002 Abrir la sesión tras verificar](TASK-002-OPEN-SESSION.md)
