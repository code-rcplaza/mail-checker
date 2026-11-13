# FOKKO-0096-SUBS · Pantalla Código

## Problema
Código verifica los 6 dígitos y abre la sesión.

## Depende de
- FOKKO-0067A-SUBS y FOKKO-0067B-SUBS, FOKKO-0091-SUBS y FOKKO-0095-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Pantalla de código, `useVerifyCode` y `useResendCode`.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/auth/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/auth`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Enviar los seis dígitos](TASK-001A-VERIFY-SUBMIT.md)
- [ ] [TASK-001B Mostrar un código inválido](TASK-001B-VERIFY-WRONG.md)
- [ ] [TASK-001C Volver a Login sin correo](TASK-001C-VERIFY-NO-EMAIL.md)
- [ ] [TASK-002 Reenviar el código](TASK-002-CODE-RESEND.md)
- [ ] [TASK-003 Verificar el código por HTTP](TASK-003-GATEWAY-VERIFY.md)
- [ ] [TASK-004 Reenviar el código por HTTP](TASK-004-GATEWAY-RESEND.md)
- [ ] [TASK-005 Conectar el hook con los componentes](TASK-005-CONTAINER-CODE.md)
