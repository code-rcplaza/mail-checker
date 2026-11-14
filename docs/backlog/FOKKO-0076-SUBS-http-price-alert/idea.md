# FOKKO-0076-SUBS · Descartar aviso de alza

## Problema
La alerta de alza se descarta con "mantener". `later` no existe en el MVP.

## Depende de
- FOKKO-0028-SUBS y FOKKO-0064-SUBS.

## Alcance
- `POST /api/subscriptions/price-alerts/decision` con `{lineKey}`.

## Fuera
- `later` → post-MVP. Requiere actualizar `docs/mvp.md` y `ARCHITECTURE.md`.

## Features afectadas
- `internal/subscriptions` (`adapters/http`)

## Dudas abiertas
- Ninguna.

## Tareas
- [x] [TASK-001 `POST /api/subscriptions/price-alerts/decision`](TASK-001-PRICE-ALERT-HANDLER.md)
