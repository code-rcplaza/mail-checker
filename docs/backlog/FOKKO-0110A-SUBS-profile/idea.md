# FOKKO-0110A-SUBS · Perfil

## Problema
Perfil muestra la cuenta y agrupa Apariencia, Tus datos y Cerrar sesión. Las filas "Cuentas de correo" y "Notificaciones" son de después del MVP.

Esta idea cubre solo: Perfil. El resto del trabajo original está en las ideas hermanas `FOKKO-0110*-SUBS`.

## Depende de
- FOKKO-0068A-SUBS y FOKKO-0068B-SUBS y FOKKO-0084A-SUBS y FOKKO-0084B-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `ProfileHeader`
- Filas y cierre de sesión

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/profile/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/profile`)

## Dudas abiertas
- Las filas post-MVP no se muestran, para no ofrecer lo que no funciona.

## Tareas
- [ ] [TASK-001 `ProfileHeader`](TASK-001-PROFILE-HEADER.md)
- [ ] [TASK-002A Filas de Perfil](TASK-002A-PROFILE-ROWS.md)
- [ ] [TASK-002B Fila Cerrar sesión](TASK-002B-PROFILE-LOGOUT-ROW.md)
- [ ] [TASK-003 Conectar el hook con los componentes](TASK-003-CONTAINER-PROFILE.md)
