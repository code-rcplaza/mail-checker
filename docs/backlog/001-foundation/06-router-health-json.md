# 001-06 · Router, health y JSON
Feature: `internal/shared/httpserver`

Como desarrollador quiero un router con health y helpers de JSON para que cada
feature solo registre sus rutas.

## Criterios
- Dado el router, cuando se hace `GET /api/health`, entonces 200,
  `Content-Type: application/json` y cuerpo `{"status":"ok"}`.
- Dado el router, cuando se pide una ruta inexistente bajo `/api`, entonces 404 con
  cuerpo `{"error":{"key":"common.route_not_found"}}`.
- Dado un cuerpo JSON válido, cuando se decodifica, entonces se obtiene el valor.
- Dado un cuerpo JSON mal formado, con campos desconocidos o de más de 1 MB, cuando
  se decodifica, entonces el error es `ErrInvalid` con la llave `common.invalid_json`.
- Dado un valor y un status, cuando se escribe la respuesta, entonces queda ese
  status y `Content-Type: application/json`.

## Fuera
- Mapeo de errores a status (001-07).
- Servir el front compilado (fase 10).
