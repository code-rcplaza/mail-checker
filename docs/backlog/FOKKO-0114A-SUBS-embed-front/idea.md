# FOKKO-0114A-SUBS · Embeber el front en el binario

## Problema
Hoy el front y el back se construyen por separado. El destino es un solo binario de Go que sirve la API y el front ya compilado, para desplegar un solo contenedor.

Esta idea cubre solo: Embeber el front en el binario. El resto del trabajo original está en las ideas hermanas `FOKKO-0114*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS y FOKKO-0082A-SUBS a FOKKO-0082D-SUBS.

## Alcance
- Embeber y servir el front

## Fuera
- Imagen de Docker → FOKKO-0115-SUBS.

## Features afectadas
- `internal/shared` (`httpserver/web`)
- `frontend` (build)

## Dudas abiertas
- Para que `go build` compile sin haber construido el front, se propone un archivo de reemplazo versionado dentro de `web/dist`. ¿Se acepta?

## Tareas
- [ ] [TASK-001A Servir `index.html`](TASK-001A-EMBED-INDEX.md)
- [ ] [TASK-001B Servir archivos con hash](TASK-001B-EMBED-ASSETS.md)
- [ ] [TASK-001C La API tiene prioridad](TASK-001C-EMBED-API-FIRST.md)
