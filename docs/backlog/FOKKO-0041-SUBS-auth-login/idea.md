# FOKKO-0041-SUBS · Pedir código al entrar (login)

## Problema
Un usuario registrado entra pidiendo un código a su correo.

## Depende de
- FOKKO-0040-SUBS, FOKKO-0042-SUBS, FOKKO-0044-SUBS y FOKKO-0045-SUBS.

## Alcance
- Caso de uso `Login {email}`: valida el correo, admite la petición según el límite, busca al usuario y le envía el código.

## Fuera
- Endpoint → FOKKO-0066A-SUBS y FOKKO-0066B-SUBS.

## Features afectadas
- `internal/auth` (`usecases`)

## Dudas abiertas
- Correo desconocido: Decidido: respuesta neutra (ver `auth` en `docs/ARCHITECTURE.md`). Responde igual que con un correo registrado y no envía nada.

## Tareas
- [x] [TASK-001A Pedir el código con un correo registrado](TASK-001A-LOGIN-REQUEST.md)
- [x] [TASK-001B Rechazar un correo inválido o desconocido](TASK-001B-LOGIN-REJECT.md)
