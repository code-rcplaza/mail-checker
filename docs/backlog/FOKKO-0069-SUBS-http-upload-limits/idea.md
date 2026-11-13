# FOKKO-0069-SUBS · Límites de la subida (tamaño, PDF, timeout)

## Problema
La subida recibe archivos de usuarios: sin tope de tamaño ni de tiempo, un archivo enorme o lento puede agotar el servidor. Fail fast en el borde.

## Alcance
- Tamaño máximo, tipo de contenido y tiempo máximo de la subida.

## Fuera
- La lógica de importar → FOKKO-0070-SUBS.

## Features afectadas
- `internal/statements` (`adapters/http`)

## Dudas abiertas
- Tamaño máximo: se propone 10 MB (una cartola en PDF pesa unos pocos MB). El rechazo usa `ErrInvalid` (400) con `statements.file_too_large`, sin agregar un error 413.

## Tareas
- [ ] [TASK-001 Tamaño máximo](TASK-001-UPLOAD-SIZE-LIMIT.md)
- [ ] [TASK-002 Formato de la subida](TASK-002-UPLOAD-MULTIPART.md)
- [ ] [TASK-003 Tiempo máximo](TASK-003-UPLOAD-TIMEOUT.md)
