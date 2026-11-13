# FOKKO-0086B-SUBS · Test de contrato con las llaves del back

## Problema
El back devuelve llaves; el front las traduce con un catálogo tipado. Un test de contrato verifica que toda llave del back exista en el catálogo.

Esta idea cubre solo: Test de contrato con las llaves del back. El resto del trabajo original está en las ideas hermanas `FOKKO-0086*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS.

## Alcance
- Contrato con el back

## Fuera
- El texto de cada pantalla: cada pantalla agrega sus llaves.

## Features afectadas
- `frontend` (`lib/i18n`)

## Dudas abiertas
- El test lee las constantes `MessageKey` de los `.go` del back con una búsqueda de texto. Si cambia la forma de declararlas, el test falla fuerte.

## Tareas
- [ ] [TASK-001 Contrato con el back](TASK-001-I18N-CONTRACT-TEST.md)
