# FOKKO-0079-SUBS/TASK-002 · Copia en JSON
Feature: `internal/auth/adapters/export`

Como usuario quiero una copia completa en JSON, para guardarla o migrarla.

## Criterios
- Dada la copia, cuando se serializa, entonces los `id` son `public_id` y no hay ids internos.
- Dada la copia, cuando se serializa, entonces las fechas son unix en segundos.

## Fuera
- Nada.

## Notas
- Ninguna.
