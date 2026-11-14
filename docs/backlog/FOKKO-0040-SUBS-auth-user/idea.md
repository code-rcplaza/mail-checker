# FOKKO-0040-SUBS · Registrar usuario

## Problema
No existe la feature `auth`. Todo dato pertenece a un usuario, así que primero hace falta el usuario y poder registrarlo (nombre, correo, términos aceptados). No hay contraseñas: el único secreto es el código por correo.

## Alcance
- `auth/domain`: `User` con sus invariantes (nombre, correo normalizado, términos aceptados).
- Puerto `UserRepository` y caso de uso para crear el usuario.

## Fuera
- Enviar el código → FOKKO-0042-SUBS y FOKKO-0045-SUBS.
- Endpoint `register` → FOKKO-0066A-SUBS y FOKKO-0066B-SUBS.

## Features afectadas
- `internal/auth` (`domain`, `ports`, `usecases`)

## Dudas abiertas
- Registrar un correo que ya existe: Decidido: respuesta neutra (ver `auth` en `docs/ARCHITECTURE.md`). No se crea otro usuario ni se devuelve error.
- Largo máximo del nombre: decidido, 60 caracteres (`auth.name_too_long`).

## Tareas
- [x] [TASK-001A Aceptar un usuario válido](TASK-001A-USER-VALID.md)
- [x] [TASK-001B Rechazar nombre o correo inválido](TASK-001B-USER-NAME-EMAIL.md)
- [x] [TASK-001C Exigir los términos](TASK-001C-USER-TERMS.md)
- [x] [TASK-002A Crear un usuario](TASK-002A-CREATE-USER.md)
- [x] [TASK-002B No duplicar un correo registrado](TASK-002B-CREATE-USER-DUPLICATE.md)
