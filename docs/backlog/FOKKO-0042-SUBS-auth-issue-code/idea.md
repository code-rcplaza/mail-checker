# FOKKO-0042-SUBS · Generar el código

## Problema
El código es de 6 dígitos, vale 10 minutos y se puede intentar 5 veces. Hay que generarlo y guardarlo de forma segura.

## Alcance
- Dominio `LoginCode` y caso de uso que emite un código nuevo (reemplaza al anterior).
- Puerto `LoginCodes` (guarda, busca, cuenta intentos).

## Fuera
- Verificarlo → FOKKO-0043-SUBS.
- Enviarlo → FOKKO-0045-SUBS.

## Features afectadas
- `internal/auth` (`domain`, `ports`, `usecases`)

## Dudas abiertas
- Se propone guardar solo un hash del código (SHA-256), no el código, y compararlo en tiempo constante. Con 6 dígitos el hash no impide un ataque offline: el límite de intentos es lo que protege.

## Tareas
- [ ] [TASK-001 `LoginCode`](TASK-001-CODE-DOMAIN.md)
- [ ] [TASK-002 Emitir un código](TASK-002-ISSUE-CODE.md)
