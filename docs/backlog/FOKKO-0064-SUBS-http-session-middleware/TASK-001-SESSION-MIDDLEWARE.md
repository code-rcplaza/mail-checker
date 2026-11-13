# FOKKO-0064-SUBS/TASK-001 · Exigir sesión
Feature: `internal/auth/adapters/http`

Como usuario quiero que nadie vea mis datos sin mi sesión, para mi privacidad.

## Criterios
- Dada una petición con una cookie de sesión válida, cuando pasa por el middleware, entonces el handler recibe el `userID`.
- Dada una petición sin cookie, cuando pasa, entonces responde 401 con `common.unauthorized` y el handler no se ejecuta.
- Dada una cookie de sesión vencida o desconocida, cuando pasa, entonces responde 401 igual que sin cookie.

## Fuera
- Nada.

## Notas
- Ninguna.
