# FOKKO-0026-SUBS · Detectar un alza

## Problema
Sin correo de aviso, el alza se detecta cuando llega el cobro mayor, con la cartola siguiente (decisión 8). Cualquier aumento cuenta (decisión 17).

## Alcance
- Detectar, por línea y moneda, un cobro mayor que el anterior: monto de antes, monto de después y desde cuándo.
- Marcar como cambio de plan el salto de más de 25 %.

## Fuera
- Impacto mensual y anual → FOKKO-0027-SUBS.
- Descartar el aviso → FOKKO-0028-SUBS.

## Features afectadas
- `internal/subscriptions` (`domain`)

## Dudas abiertas
- Los Seguros CLC suben cerca de 0,5 % al mes y alertarán cada mes. Se acepta y se revisa con usuarios reales (decisión 17).

## Tareas
- [ ] [TASK-001 Detectar el alza entre dos cobros](TASK-001-DETECT-RISE.md)
- [ ] [TASK-002 El cambio de moneda no es alza](TASK-002-RISE-IGNORES-CURRENCY-CHANGE.md)
- [ ] [TASK-003 Salto de más de 25 %](TASK-003-RISE-PLAN-CHANGE.md)
- [ ] [TASK-004 El resumen trae las alzas](TASK-004-SUMMARY-RISES.md)
