# FOKKO-0099-SUBS/TASK-003 · Fecha de facturación de respaldo
Feature: `frontend/src/features/statements/domain`

Como usuario quiero escribir la fecha de mi cartola solo cuando la app no la encuentra, para completar la subida.

## Criterios
- Dado un 400 `statements.billing_date_required`, cuando responde la subida, entonces aparece el campo de fecha con el mensaje.
- Dada una fecha escrita, cuando se reenvía, entonces se manda como `billingDate` en unix segundos.
- Dado el campo de fecha, cuando está vacío, entonces el reenvío está bloqueado.

## Fuera
- Nada.

## Notas
- Función pura en `domain`: sin React ni red.
