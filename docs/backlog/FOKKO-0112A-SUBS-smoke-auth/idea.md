# FOKKO-0112A-SUBS · Smoke test: acceso

## Problema
Cada pieza se prueba sola, pero nada prueba el flujo completo: registrarse, entrar, subir una cartola y ver el resumen. Una prueba de integración con la app real armada (base temporal, envío por consola) cubre ese hueco.

Esta idea cubre solo: Smoke test: acceso. El resto del trabajo original está en las ideas hermanas `FOKKO-0112*-SUBS`.

## Depende de
- FOKKO-0082A-SUBS a FOKKO-0082D-SUBS (y el resto del back).

## Alcance
- Registro, código y sesión

## Fuera
- Pruebas de navegador (por ejemplo Playwright): exigen una dependencia nueva y se evalúan después.

## Features afectadas
- `cmd/api` (pruebas de integración)

## Dudas abiertas
- Hace falta una cartola real en PDF anonimizada para el fixture; no se inventa. Se anonimiza antes de crearla.
- Las pruebas del front quedan con una revisión manual a 390, 768 y 1280 px por pantalla.

## Tareas
- [x] [TASK-001A Entrar con código](TASK-001A-SMOKE-LOGIN.md)
- [x] [TASK-001B Cerrar sesión](TASK-001B-SMOKE-LOGOUT.md)
- [x] [TASK-001C Cinco intentos fallidos](TASK-001C-SMOKE-ATTEMPTS.md)
