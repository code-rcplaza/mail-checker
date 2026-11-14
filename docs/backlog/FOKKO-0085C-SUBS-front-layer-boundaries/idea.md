# FOKKO-0085C-SUBS · Límites entre capas del front

## Problema
El front sigue clean architecture por feature (`domain`, `application`, `infrastructure`, `presentation`). Sin una regla automática, los componentes terminan llamando a la API y el límite se pierde. Biome debe romper el lint cuando se cruza.

## Depende de
- FOKKO-0083B-SUBS.

## Alcance
- Reglas de imports restringidos en Biome, una por límite, cada una probada con un archivo que la viola.

## Fuera
- Los gateways y el Context → FOKKO-0085D-SUBS.

## Features afectadas
- `frontend`

## Dudas abiertas
- Verificar que `noRestrictedImports` de Biome admite reglas por carpeta; si no, proponer otra herramienta (requiere aprobación).
  - Resuelto: sí, con `overrides` por carpeta en `biome.json`; no hace falta otra herramienta.

## Tareas
- [x] [TASK-001A Bloquear React y la red en `domain`](TASK-001A-DOMAIN-NO-FRAMEWORK.md)
- [x] [TASK-001B Bloquear datos y router en los componentes](TASK-001B-COMPONENTS-NO-DATA.md)
- [x] [TASK-001C Limitar `infrastructure` a la composición](TASK-001C-INFRA-ONLY-IN-MAIN.md)
