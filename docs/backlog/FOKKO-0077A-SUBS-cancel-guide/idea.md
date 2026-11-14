# FOKKO-0077A-SUBS · Guía de cancelación: caso de uso y `GET`

## Problema
"Cómo cancelar" muestra el enlace al sitio del servicio y el usuario marca "por cancelar".

Esta idea cubre solo: Guía de cancelación: caso de uso y `GET`. El resto del trabajo original está en las ideas hermanas `FOKKO-0077*-SUBS`.

## Depende de
- FOKKO-0008-SUBS, FOKKO-0035-SUBS y FOKKO-0064-SUBS.

## Alcance
- Buscar el enlace de una línea
- `GET /api/subscriptions/cancel-guide`

## Fuera
- Nada.

## Features afectadas
- `internal/subscriptions` (`usecases`, `adapters/http`)

## Dudas abiertas
- Una línea de un comercio desconocido no tiene enlace: se responde `200` con `url` vacío y el front muestra los pasos genéricos.

## Tareas
- [x] [TASK-001 Buscar el enlace de una línea](TASK-001-CANCEL-GUIDE-USECASE.md)
- [x] [TASK-002 `GET /api/subscriptions/cancel-guide`](TASK-002-CANCEL-GUIDE-HANDLER.md)
