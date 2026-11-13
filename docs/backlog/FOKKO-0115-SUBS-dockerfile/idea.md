# FOKKO-0115-SUBS · Dockerfile multi-stage

## Problema
El despliegue queda independiente de la plataforma si todo vive en una imagen de Docker: una sola imagen, configuración solo por variables de entorno y SQLite en un volumen.

## Depende de
- FOKKO-0114A-SUBS a FOKKO-0114D-SUBS y FOKKO-0081-SUBS.

## Alcance
- Dockerfile multi-stage (front, back, imagen final mínima), usuario sin privilegios, volumen para los datos y chequeo de salud.

## Fuera
- Elegir proveedor de alojamiento: sigue abierto (decisión 20).

## Features afectadas
- `.` (raíz del repo)

## Dudas abiertas
- Imagen final: se propone una base mínima (distroless o alpine). Como no trae `curl`, el chequeo de salud necesita un subcomando del binario. ¿Se acepta?
- SQLite exige una sola réplica y un volumen persistente: es el único requisito sobre la plataforma.

## Tareas
- [ ] [TASK-001 Dockerfile multi-stage](TASK-001-DOCKERFILE-BUILD.md)
- [ ] [TASK-002A Proceso sin privilegios](TASK-002A-RUNTIME-NONROOT.md)
- [ ] [TASK-002B Datos en volumen](TASK-002B-RUNTIME-VOLUME.md)
- [ ] [TASK-002C Configuración por entorno](TASK-002C-RUNTIME-ENV.md)
- [ ] [TASK-003 Chequeo de salud](TASK-003-CONTAINER-HEALTHCHECK.md)
- [ ] [TASK-004 `.dockerignore`](TASK-004-DOCKERIGNORE.md)
