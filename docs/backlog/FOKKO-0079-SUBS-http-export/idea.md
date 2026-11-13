# FOKKO-0079-SUBS · Exportar copia (CSV/JSON)

## Problema
Perfil ofrece "Descargar copia primero" antes de eliminar todo, en CSV (por defecto) o JSON. Recorre datos de varias features, así que se resuelve con un adaptador puente.

## Depende de
- FOKKO-0053A-SUBS a FOKKO-0058B-SUBS.

## Alcance
- Adaptador puente que junta los datos del usuario, formatos JSON y CSV y el endpoint.

## Fuera
- Importar una copia: post-MVP.

## Features afectadas
- `internal/auth` (`adapters/export`, puente)
- `internal/auth` (`adapters/http`)

## Dudas abiertas
- Contenido: se propone JSON con todo (perfil, cartolas, cargos, decisiones, alias, manuales, overrides) y CSV con las suscripciones del resumen (una tabla). ¿Se acepta?
- Un CSV debe escaparse contra inyección de fórmulas (`=`, `+`, `-`, `@` al inicio de un campo).

## Tareas
- [ ] [TASK-001 Juntar los datos del usuario](TASK-001-EXPORT-BRIDGE.md)
- [ ] [TASK-002 Copia en JSON](TASK-002-EXPORT-JSON.md)
- [ ] [TASK-003 Copia en CSV](TASK-003-EXPORT-CSV.md)
- [ ] [TASK-004 `GET /api/export`](TASK-004-EXPORT-HANDLER.md)
