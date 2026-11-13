# FOKKO-0069-SUBS/TASK-001 · Tamaño máximo
Feature: `internal/statements/adapters/http`

Como usuario quiero un aviso claro si mi archivo es demasiado grande, para subir uno correcto.

## Criterios
- Dado un archivo de 11 MB, cuando se sube, entonces responde 400 con `statements.file_too_large` sin leer todo el cuerpo.
- Dado un archivo de 1 MB, cuando se sube, entonces pasa.

## Fuera
- Nada.

## Notas
- Ninguna.
