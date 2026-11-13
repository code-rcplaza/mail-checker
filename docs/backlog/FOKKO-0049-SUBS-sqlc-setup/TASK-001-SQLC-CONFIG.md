# FOKKO-0049-SUBS/TASK-001 · Configuración de sqlc
Feature: `backend`

Como usuario quiero que el SQL se escriba una vez y se genere el código Go, para evitar consultas duplicadas y a mano.

## Criterios
- Dado `sqlc.yaml`, cuando se corre `sqlc generate` con una consulta de ejemplo temporal, entonces genera código Go que compila.
- Dado el esquema, entonces sale de `shared/sqlite/migrations` (una sola fuente de verdad).
- Dada la configuración, entonces cada feature escribe en su propio paquete generado.

## Fuera
- Nada.

## Notas
- La consulta de ejemplo no se commitea.
