# FOKKO-0095-SUBS · Pantalla Login

## Problema
Login pide el correo y envía un código.

## Depende de
- FOKKO-0066B-SUBS y FOKKO-0094-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Formulario de login y hook `useLogin`.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/auth/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/auth`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Enviar el login](TASK-001A-LOGIN-SUBMIT.md)
- [ ] [TASK-001B Mostrar los errores del servidor](TASK-001B-LOGIN-SERVER-ERRORS.md)
- [ ] [TASK-001C Revisión visual en tres anchos](TASK-001C-LOGIN-LAYOUT.md)
- [ ] [TASK-002A Ir de Login a Registro](TASK-002A-LOGIN-TO-REGISTER.md)
- [ ] [TASK-002B Ir de Registro a Login](TASK-002B-REGISTER-TO-LOGIN.md)
- [ ] [TASK-003 Iniciar sesión por HTTP](TASK-003-GATEWAY-LOGIN.md)
- [ ] [TASK-004 Conectar el hook con los componentes](TASK-004-CONTAINER-LOGIN.md)
