# FOKKO-0080-SUBS/TASK-003 · `DELETE /api/me`
Feature: `internal/auth/adapters/http`

Como usuario quiero confirmar y eliminar todo, para cerrar mi cuenta.

## Criterios
- Dada una sesión válida, cuando se pide, entonces responde 204 y la cookie se borra.
- Dada una petición sin sesión, cuando se pide, entonces responde 401.
- Dada la misma petición repetida con la cookie vieja, cuando se pide, entonces responde 401.

## Fuera
- Nada.

## Notas
- Ninguna.
