# FOKKO-0058A-SUBS/TASK-002 · Adaptador de líneas ocultas
Feature: `internal/subscriptions/adapters/sqlite`

Como usuario quiero que mi marca de "dejar de seguir" se conserve, para no repetirla.

## Criterios
- Dada una línea oculta dos veces, cuando se guarda, entonces queda una sola.
- Dadas marcas de dos usuarios con la misma clave, cuando uno lista, entonces ve solo las suyas.
- Dada una clave con `|`, `:`, `#` y espacios, cuando se guarda y se lee, entonces vuelve idéntica.

## Fuera
- Nada.

## Notas
- Ninguna.
