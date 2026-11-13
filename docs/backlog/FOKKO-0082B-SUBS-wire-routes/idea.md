# FOKKO-0082B-SUBS · Registrar las rutas

## Problema
`main.go` es el único lugar que conoce todas las piezas: configura, abre la base, migra, crea adaptadores, registra bancos y arma las rutas. Hoy solo levanta el router base.

Esta idea cubre solo: Registrar las rutas. El resto del trabajo original está en las ideas hermanas `FOKKO-0082*-SUBS`.

## Depende de
- FOKKO-0050-SUBS a FOKKO-0081-SUBS (según lo que exista al llegar).

## Alcance
- Registrar las rutas

## Fuera
- Servir el front → FOKKO-0114A-SUBS a FOKKO-0114D-SUBS.

## Features afectadas
- `cmd/api`

## Dudas abiertas
- Esta idea se parte en las tareas listadas; si el diff pasa de ~300 líneas en una tarea, se parte más.

## Tareas
- [ ] [TASK-001 Registrar las rutas](TASK-001-WIRE-ROUTES.md)
