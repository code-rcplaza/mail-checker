# FOKKO-0009-SUBS · `extract`: texto de un PDF

## Problema
`ImportStatement` ya llama a un `FileExtractor`, pero no existe el adaptador real. Sin él no se puede leer ninguna cartola en PDF.

## Alcance
- `statements/adapters/extract`: implementa `FileExtractor` para PDF con texto.
- Un PDF mal formado devuelve error con contexto, nunca un panic.

## Fuera
- CSV, XLSX y texto plano → post-MVP.
- OCR de PDFs escaneados: un PDF sin texto ya lo rechaza `ImportStatement` (`file_empty`).

## Features afectadas
- `internal/statements` (`adapters/extract`)

## Dudas abiertas
- Dependencia: hay que aprobar una librería de lectura de PDF (la propuesta es `ledongthuc/pdf`). CLAUDE.md exige preguntar antes de agregar dependencias.
- Hace falta un PDF real anonimizado para el fixture; no se inventa uno.

## Tareas
- [ ] [TASK-001 Extraer el texto de un PDF](TASK-001-PDF-TEXT.md)
- [ ] [TASK-002 Rechazar lo que no es un PDF](TASK-002-PDF-NOT-A-PDF.md)
- [ ] [TASK-003 PDF dañado sin panic](TASK-003-PDF-CORRUPT-NO-PANIC.md)
