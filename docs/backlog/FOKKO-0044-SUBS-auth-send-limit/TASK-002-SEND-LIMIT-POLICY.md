# FOKKO-0044-SUBS/TASK-002 · Política de 3 envíos por hora
Feature: `internal/auth/domain`

Como usuario quiero que la app limite los envíos, para que nadie abuse de mi correo.

## Criterios
- Dados 2 envíos en la última hora, cuando se evalúa, entonces se permite un tercero.
- Dados 3 envíos en la última hora, cuando se evalúa, entonces no se permite el cuarto.
- Dados 3 envíos de los cuales el más antiguo tiene más de una hora, cuando se evalúa, entonces se permite otro.

## Fuera
- Nada.

## Notas
- Ninguna.
