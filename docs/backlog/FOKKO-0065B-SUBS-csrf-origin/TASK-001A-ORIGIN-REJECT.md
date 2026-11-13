# FOKKO-0065B-SUBS/TASK-001A · Rechazar otro origen
Feature: `internal/auth/adapters/http`



## Criterios
- Dada una petición `POST`, `PUT` o `DELETE` con `Origin` de otro dominio, cuando llega, entonces se rechaza con 403.
- Dada una petición con el `Origin` de la propia app, cuando llega, entonces pasa.

## Fuera
- Nada.

## Notas
- Ninguna.
