# FOKKO-0046-SUBS/TASK-002 · Abrir la sesión tras verificar
Feature: `internal/auth/usecases`

Como usuario quiero entrar al verificar el código, para usar la app.

## Criterios
- Dado un código correcto, cuando se abre, entonces devuelve el token en claro una sola vez y guarda la sesión del usuario.
- Dado un código incorrecto, cuando se abre, entonces no se crea ninguna sesión.
- Dado un usuario con otra sesión abierta, cuando se abre una nueva, entonces ambas valen (varios dispositivos).

## Fuera
- Nada.

## Notas
- Ninguna.
