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
- Resuelta: la ruta es `/fuente`; al subir, el resultado queda en la caché
  (`["statements", "uploaded", id]`) y se navega a `/analisis/<id>` (FOKKO-0100A/B-SUBS).
- Resuelta: el texto "Qué guardamos" dice lo que guarda el back de verdad: comercio, monto, fecha y
  últimos 4 dígitos de la tarjeta de cada cobro, y el nombre del archivo; el PDF no se guarda.
  Falta la revisión legal (FOKKO-0113A-SUBS).

## Tareas
- [x] [TASK-001A Elegir o soltar un PDF](TASK-001A-DROPZONE-PICK.md)
- [x] [TASK-001B Rechazar archivos inválidos](TASK-001B-DROPZONE-REJECT.md)
- [x] [TASK-002 `BankSelect`](TASK-002-BANK-SELECT.md)
- [x] [TASK-003 Fecha de facturación de respaldo](TASK-003-BILLING-DATE-FALLBACK.md)
- [x] [TASK-004 `useUploadStatement`](TASK-004-USE-UPLOAD-STATEMENT.md)
- [x] [TASK-005 Subir la cartola por HTTP](TASK-005-GATEWAY-UPLOAD.md)
- [x] [TASK-006 Conectar el hook con los componentes](TASK-006-CONTAINER-UPLOAD.md)
