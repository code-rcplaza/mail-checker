# FOKKO-0085C-SUBS/TASK-001B · Bloquear datos y router en los componentes
Feature: `frontend`



## Criterios
- Dado un archivo de `presentation/components` o `components/ui` que importa `@tanstack/*`, `application` o `infrastructure`, cuando corre `biome lint`, entonces falla.
- Dado el mismo archivo que solo importa `react` y `lib/i18n`, entonces pasa.

## Fuera
- Nada.

## Notas
- Ninguna.
