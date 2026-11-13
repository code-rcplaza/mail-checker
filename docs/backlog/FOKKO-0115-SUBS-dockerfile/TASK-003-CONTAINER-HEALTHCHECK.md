# FOKKO-0115-SUBS/TASK-003 · Chequeo de salud
Feature: `cmd/api`

Como usuario quiero que la plataforma sepa si la app está sana, para que la reinicie si no lo está.

## Criterios
- Dado el subcomando de salud, cuando la app responde en `/api/health`, entonces termina con código 0.
- Dada la app caída o la base inaccesible, cuando se corre, entonces termina con código distinto de cero antes del tiempo máximo.
- Dado el Dockerfile, entonces declara el chequeo con ese subcomando.

## Fuera
- Nada.

## Notas
- Ninguna.
