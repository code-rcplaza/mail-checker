# FOKKO-0031-SUBS/TASK-002 · Próximo cobro de una manual
Feature: `internal/subscriptions/domain`

Como usuario quiero ver cuándo me cobran lo que agregué, para tenerlo en el calendario.

## Criterios
- Dado día de cobro 20 y hoy el día 10, cuando se calcula, entonces es el 20 de este mes.
- Dado día de cobro 5 y hoy el día 10, cuando se calcula, entonces es el 5 del mes siguiente.
- Dado día 31 y un mes de 30 días, cuando se calcula, entonces cae el último día de ese mes.

## Fuera
- Nada.

## Notas
- Ninguna.
