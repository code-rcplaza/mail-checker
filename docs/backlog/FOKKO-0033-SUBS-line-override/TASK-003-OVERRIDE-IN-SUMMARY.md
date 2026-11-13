# FOKKO-0033-SUBS/TASK-003 · El override cambia lo que se muestra
Feature: `internal/subscriptions/domain`

Como usuario quiero que el total y el próximo cobro usen mi corrección, para que cuadren con mi realidad.

## Criterios
- Dada una línea con override de monto, cuando se arma el resumen, entonces su monto y el total usan el del override.
- Dada una línea con override de día, cuando se arma el resumen, entonces su próximo cobro cae en ese día.
- Dada una línea con override, cuando se arma el resumen, entonces sus cobros reales siguen intactos.

## Fuera
- Nada.

## Notas
- Ninguna.
