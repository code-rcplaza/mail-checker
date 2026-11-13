# FOKKO-0111-SUBS · Eliminar todo

## Problema
Eliminar todo borra la cuenta y los datos, con la opción de descargar una copia antes.

## Depende de
- FOKKO-0079-SUBS, FOKKO-0080-SUBS y FOKKO-0110A-SUBS y FOKKO-0110B-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `DeleteAllConfirm`, descarga de la copia y `useDeleteAll`.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/profile/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/profile`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Descargar la copia](TASK-001-EXPORT-DOWNLOAD.md)
- [ ] [TASK-002 Confirmación](TASK-002-DELETE-ALL-CONFIRM.md)
- [ ] [TASK-003A Llamar a `DELETE /api/me`](TASK-003A-DELETE-CALL.md)
- [ ] [TASK-003B Salir tras eliminar](TASK-003B-DELETE-SUCCESS.md)
- [ ] [TASK-003C Error al eliminar](TASK-003C-DELETE-FAILURE.md)
- [ ] [TASK-004 Exportar por HTTP](TASK-004-GATEWAY-EXPORT.md)
- [ ] [TASK-005 Eliminar todo por HTTP](TASK-005-GATEWAY-DELETE-ALL.md)
- [ ] [TASK-006 Conectar el hook con los componentes](TASK-006-CONTAINER-DELETE-ALL.md)
