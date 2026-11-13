# FOKKO-0114B-SUBS/TASK-001 · Rutas del front
Feature: `internal/shared/httpserver`

Como usuario quiero recargar la página en cualquier pantalla, para no ver un 404.

## Criterios
- Dada una ruta que no es de la API ni un archivo (por ejemplo `/subs`), cuando se pide con `GET`, entonces responde `index.html` con 200.
- Dada una ruta `/api/inexistente`, cuando se pide, entonces responde 404 en JSON y no `index.html`.
- Dado un archivo estático que no existe (`/assets/falta.js`), cuando se pide, entonces responde 404 y no `index.html`.

## Fuera
- Nada.

## Notas
- Ninguna.
