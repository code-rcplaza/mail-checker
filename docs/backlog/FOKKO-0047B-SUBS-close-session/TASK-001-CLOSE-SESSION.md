# FOKKO-0047B-SUBS/TASK-001 · Cerrar la sesión
Feature: `internal/auth/usecases`

Como usuario quiero cerrar mi sesión, para que nadie más use este dispositivo.

## Criterios
- Dado un token vigente, cuando se cierra, entonces ya no valida.
- Dado un token ya cerrado o desconocido, cuando se cierra, entonces no falla (idempotente).
- Dadas dos sesiones del mismo usuario, cuando se cierra una, entonces la otra sigue valiendo.

## Fuera
- Nada.

## Notas
- Ninguna.
