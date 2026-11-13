# FOKKO-0036-SUBS/TASK-001 · Reparto por servicio
Feature: `internal/subscriptions/domain`

Como usuario quiero ver cuánto de mi gasto va a cada servicio, para saber qué pesa más.

## Criterios
- Dadas dos líneas de Netflix y una de Spotify, cuando se reparte, entonces las de Netflix suman en un solo servicio.
- Dado un reparto, cuando se calcula, entonces los porcentajes suman 100 % (el resto del redondeo va al mayor).
- Dado un gasto en CLP total cero, cuando se reparte, entonces el reparto es vacío y no hay división por cero.

## Fuera
- Nada.

## Notas
- Ninguna.
