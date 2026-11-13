# FOKKO-0035-SUBS/TASK-001 · Marcar "por cancelar"
Feature: `internal/subscriptions/usecases`

Como usuario quiero marcar una suscripción que voy a cancelar, para recordarlo.

## Criterios
- Dada una línea existente, cuando se marca, entonces se guarda con la fecha de la marca (reloj inyectado).
- Dada una línea ya marcada, cuando se marca, entonces no falla ni cambia la fecha.
- Dada una línea inexistente o de otro usuario, cuando se marca, entonces devuelve `ErrNotFound`.

## Fuera
- Nada.

## Notas
- Ninguna.
