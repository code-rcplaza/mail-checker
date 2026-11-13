# FOKKO-0104-SUBS · `CycleTip` con el texto corregido

## Problema
El tip "Optimización de ciclo" del Resumen menciona "tu cuenta corriente", lo que contradice "sin acceso a tu banco" (pendiente del handoff).

## Depende de
- FOKKO-0103-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `CycleTip` con un texto que no prometa nada que la app no sabe.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- El texto nuevo lo define el usuario; la tarea no lo inventa. Opciones a proponer al escribirla: hablar solo de las fechas de cobro de las suscripciones.

## Tareas
- [ ] [TASK-001 `CycleTip`](TASK-001-CYCLE-TIP.md)
- [ ] [TASK-002 Conectar el hook con los componentes](TASK-002-CONTAINER-CYCLE-TIP.md)
