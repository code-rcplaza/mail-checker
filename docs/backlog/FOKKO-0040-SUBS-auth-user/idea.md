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
- Registrar un correo que ya existe: responder `auth.email_taken` revela que la cuenta existe. Se propone aceptarlo en el MVP (el registro lo exige para no duplicar) y revisarlo antes de salir a producción.
- Largo máximo del nombre: se propone 60 caracteres.

## Tareas
- [ ] [TASK-001A Aceptar un usuario válido](TASK-001A-USER-VALID.md)
- [ ] [TASK-001B Rechazar nombre o correo inválido](TASK-001B-USER-NAME-EMAIL.md)
- [ ] [TASK-001C Exigir los términos](TASK-001C-USER-TERMS.md)
- [ ] [TASK-002A Crear un usuario](TASK-002A-CREATE-USER.md)
- [ ] [TASK-002B Rechazar un correo repetido](TASK-002B-CREATE-USER-DUPLICATE.md)
