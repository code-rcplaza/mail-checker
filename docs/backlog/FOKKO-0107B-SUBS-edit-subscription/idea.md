# FOKKO-0107B-SUBS · Editar una suscripción

## Problema
Editar permite nombrar una línea y corregir su monto y día; Agregar crea una suscripción manual. Un solo `SubscriptionForm` con su esquema zod.

Esta idea cubre solo: Editar una suscripción. El resto del trabajo original está en las ideas hermanas `FOKKO-0107*-SUBS`.

## Depende de
- FOKKO-0074A-SUBS y FOKKO-0074B-SUBS, FOKKO-0075-SUBS y FOKKO-0087-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Editar una línea detectada
- Editar o quitar una manual

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- El handoff deja abiertos categoría y moneda distinta de CLP: se propone no pedir categoría (YAGNI) y permitir CLP o USD.

## Tareas
- [ ] [TASK-001 Editar una línea detectada](TASK-001-EDIT-LINE-FORM.md)
- [ ] [TASK-002 Editar o quitar una manual](TASK-002-EDIT-MANUAL-FORM.md)
- [ ] [TASK-003 Editar por HTTP](TASK-003-GATEWAY-EDIT.md)
- [ ] [TASK-004 Conectar el hook con los componentes](TASK-004-CONTAINER-EDIT.md)
