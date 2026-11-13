# FOKKO-0068B-SUBS/TASK-001 · `POST /api/auth/logout`
Feature: `internal/auth/adapters/http`

Como usuario quiero cerrar mi sesión, para que nadie use este dispositivo.

## Criterios
- Dada una sesión válida, cuando se cierra, entonces responde 204, la sesión deja de valer y la cookie se borra.
- Dada una cookie ya inválida, cuando se cierra, entonces responde 204 igual.

## Fuera
- Nada.

## Notas
- Ninguna.
