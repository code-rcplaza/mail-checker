# FOKKO-0101-SUBS · Error de lectura

## Problema
Error muestra por qué falló la subida y cómo seguir. Es el 400 de `POST /api/statements` con su llave i18n.

## Depende de
- FOKKO-0070-SUBS y FOKKO-0099-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- Pantalla de error con mensaje por llave y reintento.

## Fuera
- Nada.

## Capas
clean architecture por feature (`features/statements/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/statements`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001 Pantalla de error](TASK-001-UPLOAD-ERROR-SCREEN.md)
- [ ] [TASK-002A Mensajes de banco o cartola repetida](TASK-002A-ERRMSG-BANK-IMPORTED.md)
- [ ] [TASK-002B Mensajes de archivo](TASK-002B-ERRMSG-FILE.md)
- [ ] [TASK-002C Mensaje genérico](TASK-002C-ERRMSG-UNKNOWN.md)
- [ ] [TASK-003 Conectar el hook con los componentes](TASK-003-CONTAINER-UPLOAD-ERROR.md)
