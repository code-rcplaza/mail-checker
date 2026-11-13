# FOKKO-0001-SUBS/TASK-008 · `main.go`, `slog` y apagado ordenado
Feature: `cmd/api`

Como usuario quiero que la app arranque con todo el armazón conectado y se apague
limpio para no dejar la base a medio escribir.

## Criterios
- Dada una config válida con base temporal, cuando se llama `run(ctx, ...)`, entonces
  la base queda migrada y `GET /api/health` responde 200.
- Dado `run` corriendo, cuando se cancela el contexto, entonces `run` devuelve `nil`
  en menos de 5 s y la base queda cerrada.
- Dada una config inválida, cuando se llama `run`, entonces devuelve error sin abrir
  la base ni el servidor.
- Dado el arranque, cuando se revisa el log, entonces hay una línea `slog` con el
  puerto y la ruta de la base, y ninguna con secretos.

## Fuera
- Rutas de las features (fase 5).

## Notas
- `main` solo llama `run(ctx, os.Getenv, os.Stdout)` y sale con código 1 si hay
  error; toda la lógica testeable vive en `run`.
- `ctx` se cancela con SIGINT/SIGTERM (`signal.NotifyContext`).
