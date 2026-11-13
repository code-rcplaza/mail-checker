# FOKKO-0107A-SUBS · Agregar una suscripción manual

## Problema
Editar permite nombrar una línea y corregir su monto y día; Agregar crea una suscripción manual. Un solo `SubscriptionForm` con su esquema zod.

Esta idea cubre solo: Agregar una suscripción manual. El resto del trabajo original está en las ideas hermanas `FOKKO-0107*-SUBS`.

## Depende de
- FOKKO-0074A-SUBS y FOKKO-0074B-SUBS, FOKKO-0075-SUBS y FOKKO-0087-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Esquema del formulario
- Agregar una manual

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/subscriptions/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/subscriptions`)

## Dudas abiertas
- El handoff deja abiertos categoría y moneda distinta de CLP: se propone no pedir categoría (YAGNI) y permitir CLP o USD.

## Tareas
- [ ] [TASK-001A Errores de los campos del formulario](TASK-001A-FORM-FIELDS.md)
- [ ] [TASK-001B Tipo con `z.infer`](TASK-001B-FORM-TYPE.md)
- [ ] [TASK-002 Agregar una manual](TASK-002-ADD-MANUAL-FORM.md)
- [ ] [TASK-003 Crear la manual por HTTP](TASK-003-GATEWAY-ADD-MANUAL.md)
- [ ] [TASK-004 Conectar el hook con los componentes](TASK-004-CONTAINER-ADD-MANUAL.md)
