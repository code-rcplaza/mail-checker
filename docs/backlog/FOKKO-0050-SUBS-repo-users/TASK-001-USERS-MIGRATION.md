# FOKKO-0050-SUBS/TASK-001 · Tabla `users`
Feature: `internal/shared/sqlite/migrations`

Como usuario quiero que mi cuenta quede guardada, para volver a entrar.

## Criterios
- Dada la migración, cuando se aplica y se revierte en una base temporal, entonces crea y borra `users` sin errores.
- Dada la tabla, entonces tiene `id`, `public_id` único, `name`, `email` único sin distinguir mayúsculas, `accepted_terms_at` y `created_at` (unix en segundos), todo `NOT NULL`.
- Dado un correo repetido con otra capitalización, cuando se inserta, entonces la base lo rechaza.

## Fuera
- Nada.

## Notas
- Ninguna.
