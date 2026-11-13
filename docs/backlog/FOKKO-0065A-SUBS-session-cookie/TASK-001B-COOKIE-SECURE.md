# FOKKO-0065A-SUBS/TASK-001B · Marcar `Secure` según el modo
Feature: `internal/auth/adapters/http`



## Criterios
- Dado el modo producción, cuando se fija la cookie, entonces lleva `Secure`.
- Dado el modo desarrollo (http), cuando se fija la cookie, entonces no lleva `Secure`.

## Fuera
- Nada.

## Notas
- Ninguna.
