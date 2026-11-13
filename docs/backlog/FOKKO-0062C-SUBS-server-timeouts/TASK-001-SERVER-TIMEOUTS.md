# FOKKO-0062C-SUBS/TASK-001 · Timeouts del servidor
Feature: `cmd/api`

Como usuario quiero que una conexión colgada no bloquee el servidor, para que siga respondiendo a todos.

## Criterios
- Dado el servidor, entonces tiene tiempos máximos de lectura de cabeceras, lectura, escritura e inactividad.
- Dado un cliente que no termina de enviar la cabecera, cuando pasa el tiempo máximo, entonces se corta la conexión.

## Fuera
- Nada.

## Notas
- Ninguna.
