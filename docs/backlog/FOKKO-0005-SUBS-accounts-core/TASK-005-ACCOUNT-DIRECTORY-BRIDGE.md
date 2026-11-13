# FOKKO-0005-SUBS/TASK-005 · Puente `AccountDirectory`
Feature: `internal/emailsync/adapters/accounts`

Como desarrollador quiero que `emailsync` obtenga las credenciales y registre cada
sincronización a través de `accounts`, sin que una feature importe a la otra fuera del
adaptador.

## Criterios
- Dadas cuentas en `accounts`, cuando se llama `Credentials` del puente, entonces
  devuelve `emailsync` `Credential` con id, `PublicID`, correo y contraseña, en el mismo
  orden.
- Dado un error de `accounts.Credentials` (ej. `password_unreadable`), cuando se llama
  el puente, entonces se devuelve conservando su tipo y llave.
- Dados id, hora y llave, cuando se llama `RecordSync` del puente, entonces
  `accounts.RecordSync` recibe esos valores.

## Fuera
- `AccountDirectory` en el armado de `main.go` → fase 5.

## Notas
- Depende de FOKKO-0003-SUBS/TASK-002 (`Credential`) y FOKKO-0003-SUBS/TASK-003 (puerto `AccountDirectory`).
- Se prueba con los casos de uso reales de `accounts` y dobles de sus puertos.
