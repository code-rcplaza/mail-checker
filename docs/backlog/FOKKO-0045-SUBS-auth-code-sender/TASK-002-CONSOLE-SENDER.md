# FOKKO-0045-SUBS/TASK-002 · Envío por consola para desarrollo
Feature: `internal/auth/adapters/mailer`

Como usuario quiero probar el acceso en mi máquina sin correo real, para desarrollar rápido.

## Criterios
- Dado el adaptador de consola con un `io.Writer`, cuando se envía, entonces escribe el correo y el código en ese writer.
- Dado el adaptador de consola, cuando se envía, entonces no usa `slog` (los logs nunca llevan códigos ni correos).
- Dado el modo producción, cuando `main.go` arma el adaptador, entonces se rechaza el de consola (se prueba en FOKKO-0082A-SUBS a FOKKO-0082D-SUBS).

## Fuera
- Nada.

## Notas
- Ninguna.
