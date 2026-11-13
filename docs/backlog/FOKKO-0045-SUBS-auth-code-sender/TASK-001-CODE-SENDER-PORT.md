# FOKKO-0045-SUBS/TASK-001 · Puerto de envío
Feature: `internal/auth/ports`

Como usuario quiero recibir mi código por correo, para entrar sin contraseña.

## Criterios
- Dado el puerto, cuando se envía a un correo con un código, entonces el doble en memoria registra a quién y qué código.
- Dado un fallo del envío, cuando se envía, entonces el caso de uso devuelve el error con contexto y sigue sin dar sesión.

## Fuera
- Nada.

## Notas
- Ninguna.
