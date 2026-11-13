# FOKKO-0113B-SUBS/TASK-001 · No se guarda el texto de la cartola
Feature: `internal/statements/usecases`

Como usuario quiero que la app olvide el contenido de mi cartola, para que no quede en ninguna parte.

## Criterios
- Dada una cartola con líneas que no son cobros y contienen una marca única, cuando se importa, entonces esa marca no aparece en ninguna tabla de la base.
- Dada una línea que no se pudo leer, cuando se importa, entonces solo queda su número y no su texto.
- Dado el cargo guardado, cuando se lee, entonces tiene descriptor, monto, fecha, tarjeta y últimos 4 dígitos y nada más del texto original.

## Fuera
- Nada.

## Notas
- Ninguna.
