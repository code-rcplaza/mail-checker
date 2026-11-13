# FOKKO-0102-SUBS · Revisión (con selector de cadencia)

## Problema
Revisión pide al usuario decidir sobre cada "posible": es una suscripción (con cadencia) o no lo es. Se puede deshacer.

## Depende de
- FOKKO-0073A-SUBS y FOKKO-0073B-SUBS, FOKKO-0089A-SUBS y FOKKO-0089B-SUBS y FOKKO-0090-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `ReviewCard`, `ReviewProgressRing` y los hooks de decisión.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- El selector de cadencia (mensual o anual) no está en el diseño del handoff: se agrega con `SegmentedControl`.

## Tareas
- [ ] [TASK-001 Lista de posibles](TASK-001-REVIEW-LIST.md)
- [ ] [TASK-002 Confirmar con cadencia](TASK-002-REVIEW-CONFIRM.md)
- [ ] [TASK-003A Descartar una posible](TASK-003A-REVIEW-DISMISS.md)
- [ ] [TASK-003B Deshacer una decisión](TASK-003B-REVIEW-UNDO.md)
- [ ] [TASK-004 Progreso de la revisión](TASK-004-REVIEW-PROGRESS.md)
- [ ] [TASK-005 Leer las pendientes de revisión](TASK-005-HOOK-REVIEW-LIST.md)
- [ ] [TASK-006 Leer las pendientes por HTTP](TASK-006-GATEWAY-REVIEW-READ.md)
- [ ] [TASK-007 Enviar una decisión por HTTP](TASK-007-GATEWAY-REVIEW-DECIDE.md)
- [ ] [TASK-008 Conectar el hook con los componentes](TASK-008-CONTAINER-REVIEW.md)
