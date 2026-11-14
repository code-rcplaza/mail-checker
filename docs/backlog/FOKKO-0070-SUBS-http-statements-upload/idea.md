# FOKKO-0070-SUBS · `POST /api/statements`

## Problema
`POST /api/statements` sube una cartola y responde archivo y meses leídos. Es el centro del flujo de primera vez.

## Depende de
- FOKKO-0010-SUBS, FOKKO-0018A-SUBS y FOKKO-0018B-SUBS, FOKKO-0020A-SUBS y FOKKO-0020B-SUBS, FOKKO-0064-SUBS y FOKKO-0069-SUBS.

## Alcance
- Handler multipart con `file`, `billingDate` y `bank` opcionales.

## Fuera
- `GET` y `DELETE` de cartolas: post-MVP.

## Features afectadas
- `internal/statements` (`adapters/http`)

## Dudas abiertas
- `billingDate` viaja como unix en segundos, igual que el resto de la API.

## Tareas
- [x] [TASK-001 Subir una cartola](TASK-001-STATEMENT-UPLOAD-HANDLER.md)
- [x] [TASK-002A Errores 400 de la subida](TASK-002A-UPLOAD-ERR-400.md)
- [x] [TASK-002B Error 409 de cartola repetida](TASK-002B-UPLOAD-ERR-409.md)
- [x] [TASK-003A Campo `billingDate`](TASK-003A-BILLING-DATE-FIELD.md)
- [x] [TASK-003B Campo `bank`](TASK-003B-BANK-FIELD.md)
