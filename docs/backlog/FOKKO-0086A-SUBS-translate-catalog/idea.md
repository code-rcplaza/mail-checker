# FOKKO-0086A-SUBS · Traducción tipada con catálogo en español

## Problema
El back devuelve llaves; el front las traduce con un catálogo tipado. Un test de contrato verifica que toda llave del back exista en el catálogo.

Esta idea cubre solo: Traducción tipada con catálogo en español. El resto del trabajo original está en las ideas hermanas `FOKKO-0086*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS.

## Alcance
- `t(key, params)` tipada
- Catálogo en español

## Fuera
- El texto de cada pantalla: cada pantalla agrega sus llaves.

## Features afectadas
- `frontend` (`lib/i18n`)

## Dudas abiertas
- El test lee las constantes `MessageKey` de los `.go` del back con una búsqueda de texto. Si cambia la forma de declararlas, el test falla fuerte.

## Tareas
- [x] [TASK-001A Traducir una llave](TASK-001A-T-LOOKUP.md)
- [x] [TASK-001B Mensaje genérico para llaves desconocidas](TASK-001B-T-FALLBACK.md)
- [x] [TASK-002 Catálogo en español](TASK-002-CATALOG-ES.md)
