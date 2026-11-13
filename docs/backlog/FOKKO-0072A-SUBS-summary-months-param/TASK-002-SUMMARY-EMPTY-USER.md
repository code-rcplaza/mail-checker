# FOKKO-0072A-SUBS/TASK-002 · Usuario sin datos
Feature: `internal/subscriptions/adapters/http`

Como usuario quiero ver un resumen vacío si no cargué nada, para empezar desde cero.

## Criterios
- Dado un usuario sin cartolas, cuando se pide el resumen, entonces responde 200 con las listas vacías (no `null`) y total en cero.
- Dada una petición sin sesión, cuando se pide, entonces responde 401.

## Fuera
- Nada.

## Notas
- Ninguna.
