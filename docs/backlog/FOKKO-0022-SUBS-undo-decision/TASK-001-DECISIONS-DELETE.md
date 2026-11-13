# FOKKO-0022-SUBS/TASK-001 · `Delete` en el puerto de decisiones
Feature: `internal/subscriptions/ports`

Como usuario quiero poder borrar una decisión guardada, para arrepentirme.

## Criterios
- Dada una decisión guardada, cuando se llama `Delete` con su clave, entonces ya no aparece en `List`.
- Dada una clave sin decisión, cuando se llama `Delete`, entonces no falla.

## Fuera
- Nada.

## Notas
- Ninguna.
