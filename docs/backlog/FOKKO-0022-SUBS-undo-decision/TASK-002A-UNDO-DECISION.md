# FOKKO-0022-SUBS/TASK-002A · Deshacer una decisión
Feature: `internal/subscriptions/usecases`



## Criterios
- Dada una línea confirmada, cuando se deshace, entonces el siguiente `GetSummary` la muestra en posibles.
- Dada una línea descartada, cuando se deshace, entonces vuelve a posibles.
- Dada una línea sin decisión, cuando se deshace, entonces no falla (idempotente).

## Fuera
- Nada.

## Notas
- Ninguna.
