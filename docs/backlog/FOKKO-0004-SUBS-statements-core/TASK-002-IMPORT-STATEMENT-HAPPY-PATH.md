# FOKKO-0004-SUBS/TASK-002 · `ImportStatement`: camino feliz con banco elegido
Feature: `internal/statements/usecases`

Como usuario quiero importar una cartola de un banco que elijo para que sus cobros
queden guardados.

## Criterios
- Dado un archivo, una fecha de facturación y un banco registrado, cuando se llama
  `ImportStatement`, entonces se extrae el texto, se parsea con el parser de ese banco y
  se llama `Create` una sola vez con la cartola y sus cobros.
- Dada una importación exitosa, entonces la `Statement` guardada tiene el `PublicID` del
  generador, `ImportedAt` del reloj, el nombre del archivo, el `FileKind` que devolvió el
  extractor, la fecha de facturación, la cantidad de cobros y las líneas malas (las del
  parser más las que no se pudieron convertir).
- Dada una importación exitosa, entonces `ImportStatement` devuelve la `Statement`
  guardada.

## Fuera
- Detección del banco → FOKKO-0004-SUBS/TASK-003. Validaciones y errores → FOKKO-0004-SUBS/TASK-004.

## Notas
- Puertos `FileExtractor`, `BankStatementParser`, `StatementRepository` en `ports`;
  dobles en memoria en los tests.
- Generador de `PublicID` y reloj inyectados.
