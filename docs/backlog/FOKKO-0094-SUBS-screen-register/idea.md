# FOKKO-0094-SUBS · Pantalla Registro

## Problema
Registro es la primera pantalla de un usuario nuevo: nombre, correo y aceptar términos. Valida en el front con zod antes de llamar al back.

## Depende de
- FOKKO-0066A-SUBS, FOKKO-0087-SUBS, FOKKO-0092-SUBS y FOKKO-0086A-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Esquema zod, formulario y hook `useRegister`.

## Fuera
- Botones de Google y Apple (el diseño los trae; el MVP no).
- Los textos legales → FOKKO-0113A-SUBS y FOKKO-0113B-SUBS.

## Capas
clean architecture por feature (`features/auth/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/auth`)

## Dudas abiertas
- Los términos y la política de privacidad deben enlazar a páginas que aún no existen (FOKKO-0113A-SUBS y FOKKO-0113B-SUBS); hasta entonces el texto va sin enlace.
- TASK-002C queda abierta: no hay navegador en la máquina para revisar a 390, 768 y 1280 px.
- Resuelta: el correo va a Código en la búsqueda (`/codigo?email=`); el back no registra la query en los logs (FOKKO-0062B-SUBS).

## Tareas
- [x] [TASK-001A Errores del nombre o el correo](TASK-001A-SCHEMA-NAME-EMAIL.md)
- [x] [TASK-001B Exigir los términos](TASK-001B-SCHEMA-TERMS.md)
- [x] [TASK-001C Tipo con `z.infer`](TASK-001C-SCHEMA-TYPE.md)
- [x] [TASK-002A Enviar el registro](TASK-002A-REGISTER-SUBMIT.md)
- [x] [TASK-002B Mostrar los errores del servidor](TASK-002B-REGISTER-SERVER-ERRORS.md)
- [ ] [TASK-002C Revisión visual en tres anchos](TASK-002C-REGISTER-LAYOUT.md)
- [x] [TASK-003 Registrar por HTTP](TASK-003-GATEWAY-REGISTER.md)
- [x] [TASK-004 Conectar el hook con los componentes](TASK-004-CONTAINER-REGISTER.md)
