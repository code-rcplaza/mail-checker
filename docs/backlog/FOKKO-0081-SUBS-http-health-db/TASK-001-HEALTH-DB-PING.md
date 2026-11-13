# FOKKO-0081-SUBS/TASK-001 · La salud comprueba la base
Feature: `internal/shared/httpserver`

Como usuario quiero que la plataforma reinicie la app si pierde la base, para que se recupere sola.

## Criterios
- Dada una base que responde, cuando se pide `/api/health`, entonces responde 200 `{status:"ok"}`.
- Dada una base cerrada o que no responde en el tiempo máximo, cuando se pide, entonces responde 503 con `common.unavailable`.
- Dada la respuesta de salud, entonces no revela detalles internos ni la ruta de la base.

## Fuera
- Nada.

## Notas
- Ninguna.
