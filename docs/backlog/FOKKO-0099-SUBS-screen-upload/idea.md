# FOKKO-0099-SUBS · Subir cartola

## Problema
Elegir fuente solo ofrece subir la cartola (PDF). El usuario puede arrastrar o elegir un archivo y, si hace falta, el banco o la fecha.

## Depende de
- FOKKO-0070-SUBS, FOKKO-0071-SUBS y FOKKO-0087-SUBS.
- FOKKO-0085C-SUBS y FOKKO-0085D-SUBS.

## Alcance
- `UploadDropzone`, `BankSelect`, campo de fecha de respaldo y `useUploadStatement`.

## Fuera
- Cargando, Análisis y Error → FOKKO-0100A-SUBS y FOKKO-0100B-SUBS y FOKKO-0101-SUBS.

## Capas
clean architecture por feature (`features/statements/{domain,application,infrastructure,presentation}`); los componentes son presentacionales y el contenedor de `screens` los conecta con los hooks.

## Features afectadas
- `frontend` (`features/statements`)

## Dudas abiertas
- Ninguna.

## Tareas
- [ ] [TASK-001A Elegir o soltar un PDF](TASK-001A-DROPZONE-PICK.md)
- [ ] [TASK-001B Rechazar archivos inválidos](TASK-001B-DROPZONE-REJECT.md)
- [ ] [TASK-002 `BankSelect`](TASK-002-BANK-SELECT.md)
- [ ] [TASK-003 Fecha de facturación de respaldo](TASK-003-BILLING-DATE-FALLBACK.md)
- [ ] [TASK-004 `useUploadStatement`](TASK-004-USE-UPLOAD-STATEMENT.md)
- [ ] [TASK-005 Subir la cartola por HTTP](TASK-005-GATEWAY-UPLOAD.md)
- [ ] [TASK-006 Conectar el hook con los componentes](TASK-006-CONTAINER-UPLOAD.md)
