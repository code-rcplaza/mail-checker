# FOKKO-0077B-SUBS · `POST /api/subscriptions/cancellations`

## Problema
"Cómo cancelar" muestra el enlace al sitio del servicio y el usuario marca "por cancelar".

Esta idea cubre solo: `POST /api/subscriptions/cancellations`. El resto del trabajo original está en las ideas hermanas `FOKKO-0077*-SUBS`.

## Depende de
- FOKKO-0008-SUBS, FOKKO-0035-SUBS y FOKKO-0064-SUBS.

## Alcance
- `POST /api/subscriptions/cancellations`

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`usecases`, `adapters/http`)

## Dudas abiertas
- Una línea de un comercio desconocido no tiene enlace: se responde `200` con `url` vacío y el front muestra los pasos genéricos.

## Tareas
- [ ] [TASK-001 `POST /api/subscriptions/cancellations`](TASK-001-CANCELLATION-HANDLER.md)
