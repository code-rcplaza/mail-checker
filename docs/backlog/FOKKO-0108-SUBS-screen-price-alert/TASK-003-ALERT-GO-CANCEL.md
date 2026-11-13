# FOKKO-0108-SUBS/TASK-003 · Ir a cancelar
Feature: `frontend/src/features/subscriptions/presentation/screens`

Como usuario quiero ir a cancelar si el precio no me sirve, para dejar de pagarlo.

## Criterios
- Dado "Quiero cancelar", cuando se toca, entonces navega a Cómo cancelar de esa línea.
- Dada la alerta, cuando se renderiza, entonces no ofrece "después" (no existe en el MVP).

## Fuera
- Nada.

## Notas
- Contenedor en `presentation/screens`: llama hooks de `application` y pasa props.
