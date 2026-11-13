# FOKKO-0041-SUBS · Pedir código al entrar (login)

## Problema
Un usuario registrado entra pidiendo un código a su correo.

## Depende de
- FOKKO-0040-SUBS y FOKKO-0042-SUBS.

## Alcance
- Caso de uso `Login {email}`: valida el correo, busca al usuario y pide el código.

## Fuera
- Endpoint → FOKKO-0066A-SUBS y FOKKO-0066B-SUBS.

## Features afectadas
- `internal/auth` (`usecases`)

## Dudas abiertas
- Correo desconocido: responder error revela qué correos existen. Se propone responder `auth.user_not_found` en el MVP, igual que con el registro, y revisarlo antes de producción.

## Tareas
- [ ] [TASK-001A Pedir el código con un correo registrado](TASK-001A-LOGIN-REQUEST.md)
- [ ] [TASK-001B Rechazar un correo inválido o desconocido](TASK-001B-LOGIN-REJECT.md)
