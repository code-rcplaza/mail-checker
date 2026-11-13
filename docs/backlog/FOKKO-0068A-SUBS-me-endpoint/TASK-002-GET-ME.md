# FOKKO-0068A-SUBS/TASK-002 · `GET /api/me`
Feature: `internal/auth/adapters/http`

Como usuario quiero ver mi nombre y correo en Perfil, para confirmar con qué cuenta estoy.

## Criterios
- Dada una sesión válida, cuando se pide, entonces responde `{id, name, email, initials}` con `id` igual al `public_id`.
- Dada una sesión válida, cuando se pide, entonces el cuerpo no contiene el `id` interno.
- Dada una petición sin sesión, cuando se pide, entonces responde 401.

## Fuera
- Nada.

## Notas
- Ninguna.
