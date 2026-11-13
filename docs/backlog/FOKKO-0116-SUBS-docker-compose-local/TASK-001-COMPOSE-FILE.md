# FOKKO-0116-SUBS/TASK-001 · `docker-compose.yml`
Feature: `.`

Como usuario quiero levantar la app completa con un comando, para probarla como en producción.

## Criterios
- Dado `docker compose up`, cuando arranca, entonces la app responde en el puerto publicado y el volumen guarda la base.
- Dado `docker compose down` y luego `up`, cuando se repite, entonces los datos siguen.
- Dado el servicio, entonces declara el chequeo de salud de la imagen.

## Fuera
- Nada.

## Notas
- Ninguna.
