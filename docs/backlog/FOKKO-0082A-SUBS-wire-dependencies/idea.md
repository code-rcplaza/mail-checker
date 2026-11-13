# FOKKO-0082A-SUBS · Armar las dependencias

## Problema
`main.go` es el único lugar que conoce todas las piezas: configura, abre la base, migra, crea adaptadores, registra bancos y arma las rutas. Hoy solo levanta el router base.

Esta idea cubre solo: Armar las dependencias. El resto del trabajo original está en las ideas hermanas `FOKKO-0082*-SUBS`.

## Depende de
- FOKKO-0050-SUBS a FOKKO-0081-SUBS (según lo que exista al llegar).

## Alcance
- Crear los adaptadores y casos de uso

## Fuera
- Servir el front → FOKKO-0114A-SUBS a FOKKO-0114D-SUBS.

## Features afectadas
- `cmd/api`

## Dudas abiertas
- Esta idea se parte en las tareas listadas; si el diff pasa de ~300 líneas en una tarea, se parte más.

## Tareas
- [ ] [TASK-001A Armar el núcleo de la app](TASK-001A-WIRE-CORE.md)
- [ ] [TASK-001B Registrar los bancos una sola vez](TASK-001B-WIRE-BANKS.md)
- [ ] [TASK-001C Elegir el envío de códigos por modo](TASK-001C-WIRE-SENDER.md)
