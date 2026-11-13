# FOKKO-0114C-SUBS · Cabeceras de caché

## Problema
Hoy el front y el back se construyen por separado. El destino es un solo binario de Go que sirve la API y el front ya compilado, para desplegar un solo contenedor.

Esta idea cubre solo: Cabeceras de caché. El resto del trabajo original está en las ideas hermanas `FOKKO-0114*-SUBS`.

## Depende de
- FOKKO-0083A-SUBS a FOKKO-0083C-SUBS y FOKKO-0082A-SUBS a FOKKO-0082D-SUBS.

## Alcance
- Cabeceras de caché

## Fuera
- Imagen de Docker → FOKKO-0115-SUBS.

## Features afectadas
- `internal/shared` (`httpserver/web`)
- `frontend` (build)

## Dudas abiertas
- Para que `go build` compile sin haber construido el front, se propone un archivo de reemplazo versionado dentro de `web/dist`. ¿Se acepta?

## Tareas
- [ ] [TASK-001 Cabeceras de caché](TASK-001-STATIC-CACHE.md)
