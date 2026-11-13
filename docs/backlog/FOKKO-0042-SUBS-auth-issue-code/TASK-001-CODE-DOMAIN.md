# FOKKO-0042-SUBS/TASK-001 · `LoginCode`
Feature: `internal/auth/domain`

Como usuario quiero un código corto que caduque, para que no sirva si alguien lo ve más tarde.

## Criterios
- Dada una fuente de números aleatoria inyectada, cuando se genera, entonces el código tiene exactamente 6 dígitos, con ceros a la izquierda si hace falta.
- Dado el reloj en un instante T, cuando se genera, entonces vence en T + 10 minutos y lleva 0 intentos.
- Dada una fuente aleatoria que falla, cuando se genera, entonces devuelve el error con contexto.

## Fuera
- Nada.

## Notas
- Las constantes (6 dígitos, 10 minutos, 5 intentos) son constantes con nombre en `auth/domain`.
