# 001-02 · Configuración con fail fast
Feature: `cmd/api/config.go`

Como usuario quiero que la app no arranque si la configuración está mal para
enterarme al inicio y no en el primer uso.

## Criterios
- Dado que no hay variables definidas, cuando se carga la config, entonces
  `Port = 8080` y `DBPath = "data/mail-checker.db"`.
- Dado `MAILCHECKER_PORT=9090` y `MAILCHECKER_DB_PATH=/tmp/x.db`, cuando se carga,
  entonces se usan esos valores.
- Dado `MAILCHECKER_PORT=abc`, cuando se carga, entonces error que nombra
  `MAILCHECKER_PORT`.
- Dado `MAILCHECKER_PORT=0` o `70000`, cuando se carga, entonces error (fuera de 1–65535).
- Dado `MAILCHECKER_DB_PATH` con solo espacios, cuando se carga, entonces error que
  nombra `MAILCHECKER_DB_PATH`.

## Fuera
- `MAILCHECKER_KEY` (llave de cifrado, fase 4).

## Notas
- La carga recibe una función `getenv func(string) string` para testear sin tocar
  el entorno real.
- Valores por defecto propuestos: puerto 8080, base en `data/` (ya ignorada en git).
